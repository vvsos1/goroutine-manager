FROM golang:1.23 AS builder

WORKDIR /app

# tls 인증서 신뢰 문제로 인해 아래 설정 추가
ENV GOPROXY=direct
ENV GOINSECURE=golang.org/*
RUN git config --global http.sslVerify false

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 go build -o /app/server ./cmd/server


FROM scratch

COPY --from=builder /app/server /server

# 기본 환경변수 및 포트
ENV VALKEY_ADDR=localhost:6379
ENV PORT=3000

EXPOSE 3000

ENTRYPOINT ["/server"]