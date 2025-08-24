FROM golang:1.23.12-alpine as builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN go build -o server ./cmd/server

FROM scratch

COPY --from=builder /app/server /server

# 기본 환경변수 및 포트
ENV VALKEY_ADDR=localhost:6379
ENV PORT=3000

EXPOSE 3000

ENTRYPOINT ["/server"]