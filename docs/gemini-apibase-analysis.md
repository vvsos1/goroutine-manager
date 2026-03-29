# Gemini Provider `apiBase` 미반영 버그 분석 보고서

> 대상 저장소: [continuedev/continue](https://github.com/continuedev/continue)  
> 분석 기준 커밋: `d971363a7e2711a056a2daafb4a788140bbf32a4`  
> 분석 일시: 2026-03-29

---

## 1. 요약

**결론: 버그 확인됨.**

Continue GoLand 플러그인에서 `provider: gemini`로 설정 시, `config.json`에 지정한 `apiBase` 값이 **채팅(chat) 요청의 실제 HTTP URL에 반영되지 않습니다.**  
임베딩(embed) 요청은 `apiBase`를 올바르게 사용하지만, 채팅 스트림 요청은 항상 Google의 기본 엔드포인트(`https://generativelanguage.googleapis.com/v1beta/`)로 고정됩니다.

---

## 2. 재현 조건

`~/.continue/config.json`에 다음과 같이 설정:

```json
{
  "models": [
    {
      "title": "Gemini Custom",
      "provider": "gemini",
      "model": "gemini-2.5-flash",
      "apiKey": "AIza...",
      "apiBase": "https://my-proxy.example.com/v1beta/"
    }
  ]
}
```

- **기대 동작**: 채팅 요청이 `https://my-proxy.example.com/v1beta/models/gemini-2.5-flash:streamGenerateContent` 로 전송
- **실제 동작**: 채팅 요청이 `https://generativelanguage.googleapis.com/v1beta/models/gemini-2.5-flash:streamGenerateContent` 로 전송 (커스텀 `apiBase` 무시)

---

## 3. 코드 경로 분석

### 3-1. 요청 라우팅 흐름

```
사용자 설정 (config.json)
  └─► core/llm/llms/Gemini.ts  ── Gemini 클래스
        ├─ useOpenAIAdapterFor = ["chat","embed","list","rerank","streamChat","streamFim"]
        └─ BaseLLM.createOpenAiAdapter()
             └─► constructLlmApi({ provider:"gemini", apiBase: this.apiBase })
                   └─► new GeminiApi(config)           ← packages/openai-adapters/src/apis/Gemini.ts
                         └─► chatCompletionStream()
                               └─► generateStream(this.genAI, ...)
                                     └─► genAI.models.generateContentStream(...)  ← @google/genai SDK
```

### 3-2. 버그 위치: `packages/openai-adapters/src/apis/Gemini.ts`

#### 현재 코드 (버그)

```typescript
// packages/openai-adapters/src/apis/Gemini.ts

export class GeminiApi implements BaseLlmApi {
  apiBase: string = "https://generativelanguage.googleapis.com/v1beta/";
  private genAI: GoogleGenAI;

  constructor(protected config: GeminiConfig) {
    this.apiBase = config.apiBase ?? this.apiBase;  // ✅ this.apiBase 에는 올바르게 저장
    this.genAI = withNativeFetch(
      () =>
        new GoogleGenAI({
          apiKey: this.config.apiKey,
          // ❌ httpOptions.baseUrl 미전달 → SDK가 항상 기본 Google 엔드포인트 사용
        }),
    );
  }

  async *chatCompletionStream(
    body: ChatCompletionCreateParamsStreaming,
    _signal: AbortSignal,
  ): AsyncGenerator<ChatCompletionChunk> {
    const convertedBody = this._convertBody(
      body,
      this.apiBase.includes("/v1/"),  // this.apiBase 는 정상적으로 커스텀 값
      true,
    );
    const response = await this.generateStream(
      this.genAI,  // ❌ 커스텀 baseUrl 없이 생성된 genAI 인스턴스
      body.model,
      convertedBody,
    );
    yield* this.processStreamResponse(response, body.model);
  }

  async embed(body: EmbeddingCreateParams): Promise<CreateEmbeddingResponse> {
    const resp = await customFetch(this.config.requestOptions)(
      new URL(`${body.model}:batchEmbedContents`, this.apiBase),  // ✅ 임베딩은 apiBase 올바르게 사용
      // ...
    );
  }
}
```

#### 핵심 문제

| 요청 유형 | `apiBase` 반영 여부 | 비고 |
|-----------|-------------------|------|
| `chatCompletionStream` | ❌ **미반영** | `GoogleGenAI` SDK 기본 엔드포인트 고정 |
| `chatCompletionNonStream` | ❌ **미반영** | `chatCompletionStream` 경유 |
| `embed` | ✅ 반영 | `customFetch` + `this.apiBase`로 직접 URL 구성 |

### 3-3. `@google/genai` SDK URL 구성 메커니즘

`GoogleGenAI` 생성자는 `httpOptions.baseUrl`과 `httpOptions.apiVersion`으로 요청 URL을 구성합니다:

```
최종 URL = ${baseUrl (trailing slash 제거)} / ${apiVersion} / ${path}
기본값:   https://generativelanguage.googleapis.com  /  v1beta  /  models/{model}:generateContent
```

`httpOptions`를 전달하지 않으면 항상 기본값(`https://generativelanguage.googleapis.com/v1beta/`)을 사용합니다.

소스 근거: `googleapis/js-genai` → [`src/_api_client.ts`](https://github.com/googleapis/js-genai/blob/585c17747d903d599f8f24d204ac245661831698/src/_api_client.ts#L242)

```typescript
// _api_client.ts, line 237-242
} else {
  // Gemini API
  if (!this.clientOptions.apiKey) {
    console.warn('API key should be set when using the Gemini API.');
  }
  initHttpOptions.apiVersion =
    this.clientOptions.apiVersion ?? GOOGLE_AI_API_DEFAULT_VERSION;  // "v1beta"
  initHttpOptions.baseUrl = `https://generativelanguage.googleapis.com/`;  // ← 하드코딩
}
```

---

## 4. 영향 범위

- **영향 받는 기능**: gemini provider를 사용하는 모든 **채팅** 기능 (대화, 코드 설명, 요약 등)
- **영향 받지 않는 기능**: 임베딩 (임베딩은 `customFetch`로 직접 URL 구성)
- **영향 받는 사용자**: 
  - 기업 방화벽 내부 Gemini API 프록시 사용자
  - Google의 Gemini API 호환 커스텀 엔드포인트 사용자  
  - AI 게이트웨이/미들웨어를 통해 Gemini를 사용하는 사용자
- **영향 없는 사용자**: `apiBase`를 설정하지 않고 기본 Google 엔드포인트를 사용하는 사용자

---

## 5. 수정 방향 및 최소 패치

### 수정 원칙

`GoogleGenAI` 인스턴스 생성 시 `httpOptions.baseUrl`에 `this.apiBase`를 전달합니다.  
Continue의 `apiBase`는 이미 API 버전 경로(`/v1beta/`)를 포함하므로, SDK가 자체적으로 `/v1beta`를 덧붙이지 않도록 `apiVersion: ""`을 함께 설정합니다.

### 패치 파일

> 전체 diff: [`docs/patches/0001-fix-gemini-apibase.patch`](./patches/0001-fix-gemini-apibase.patch)

```diff
--- a/packages/openai-adapters/src/apis/Gemini.ts
+++ b/packages/openai-adapters/src/apis/Gemini.ts
@@ -71,9 +71,15 @@ export class GeminiApi implements BaseLlmApi {
   constructor(protected config: GeminiConfig) {
     this.apiBase = config.apiBase ?? this.apiBase;
+    // Create GoogleGenAI with native fetch to avoid pollution
+    // from Vercel AI SDK packages that can break stream handling
     this.genAI = withNativeFetch(
       () =>
         new GoogleGenAI({
           apiKey: this.config.apiKey,
+          // Pass apiBase so custom endpoints are respected.
+          // apiVersion is set to "" because apiBase already contains the
+          // version segment (e.g. /v1beta/), preventing the SDK from
+          // appending it a second time.
+          httpOptions: {
+            baseUrl: this.apiBase,
+            apiVersion: "",
+          },
         }),
     );
   }
```

### 수정 효과 검증

| `apiBase` 설정 값 | 수정 전 요청 URL | 수정 후 요청 URL |
|---|---|---|
| 미설정 (기본값) | `.../v1beta/models/...` ✅ | `.../v1beta/models/...` ✅ |
| `https://proxy.example.com/v1beta/` | `googleapis.com/v1beta/models/...` ❌ | `proxy.example.com/v1beta/models/...` ✅ |
| `https://proxy.example.com/v1/` | `googleapis.com/v1beta/models/...` ❌ | `proxy.example.com/v1/models/...` ✅ |

---

## 6. 정적 검증 (단위 테스트)

패치 적용 후 다음 테스트를 `packages/openai-adapters/src/test/main.test.ts`에 추가하여 검증합니다:

```typescript
// packages/openai-adapters/src/test/main.test.ts 의 Configuration describe 블록에 추가

it("should configure GeminiApi with default apiBase when not specified", () => {
  const gemini = constructLlmApi({
    provider: "gemini",
    apiKey: "AIza-test",
  }) as any;

  expect(gemini.apiBase).toBe("https://generativelanguage.googleapis.com/v1beta/");
  // SDK 내부의 httpOptions.baseUrl 이 apiBase 로 설정됐는지 확인
  expect(
    gemini.genAI.apiClient.clientOptions.httpOptions?.baseUrl
  ).toBe("https://generativelanguage.googleapis.com/v1beta/");
});

it("should configure GeminiApi with custom apiBase", () => {
  const customBase = "https://my-proxy.example.com/v1beta/";
  const gemini = constructLlmApi({
    provider: "gemini",
    apiKey: "AIza-test",
    apiBase: customBase,
  }) as any;

  expect(gemini.apiBase).toBe(customBase);
  // 핵심: SDK 인스턴스도 커스텀 baseUrl 을 사용해야 함
  expect(
    gemini.genAI.apiClient.clientOptions.httpOptions?.baseUrl
  ).toBe(customBase);
  // apiVersion 이 "" 로 설정되어 이중 버전 경로 방지
  expect(
    gemini.genAI.apiClient.clientOptions.httpOptions?.apiVersion
  ).toBe("");
});

it("should configure GeminiApi with v1 apiBase", () => {
  const v1Base = "https://my-proxy.example.com/v1/";
  const gemini = constructLlmApi({
    provider: "gemini",
    apiKey: "AIza-test",
    apiBase: v1Base,
  }) as any;

  expect(gemini.apiBase).toBe(v1Base);
  expect(
    gemini.genAI.apiClient.clientOptions.httpOptions?.baseUrl
  ).toBe(v1Base);
});
```

> **테스트 실행 방법** (continuedev/continue 저장소에서):
> ```bash
> cd packages/openai-adapters
> IGNORE_API_KEY_TESTS=true npx vitest run src/test/main.test.ts
> ```

---

## 7. 관련 파일 및 라인 참조

| 파일 | 관련 내용 |
|------|----------|
| `core/llm/llms/Gemini.ts` | `useOpenAIAdapterFor = ["chat","embed",...]` → openai-adapters 경로로 라우팅 |
| `core/llm/index.ts` (BaseLLM) | `createOpenAiAdapter()`: `constructLlmApi({..., apiBase: this.apiBase})` 호출 |
| `packages/openai-adapters/src/index.ts` | `case "gemini": return new GeminiApi(config)` |
| `packages/openai-adapters/src/apis/Gemini.ts` | **버그 위치**: `GoogleGenAI` 생성자에 `httpOptions.baseUrl` 누락 |
| `googleapis/js-genai src/_api_client.ts` | SDK URL 구성 로직: `baseUrl + apiVersion + path` |
| `googleapis/js-genai src/types.ts` | `HttpOptions.baseUrl`, `HttpOptions.apiVersion` 타입 정의 |

---

## 8. 결론

- **버그 확인**: `packages/openai-adapters/src/apis/Gemini.ts` 의 `GeminiApi` 생성자에서 `GoogleGenAI` SDK 인스턴스를 생성할 때 `httpOptions.baseUrl`을 전달하지 않아, 사용자가 설정한 `apiBase`가 채팅 요청에 반영되지 않음
- **최소 수정**: 생성자에 `httpOptions: { baseUrl: this.apiBase, apiVersion: "" }` 3줄 추가
- **하위 호환**: 기본값(`apiBase` 미설정) 사용 시 동작 변화 없음
- **테스트**: 단위 테스트 3개로 정적 검증 가능
