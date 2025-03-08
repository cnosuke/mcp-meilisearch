FROM golang:1.24 as builder

WORKDIR /app

# 依存関係のコピーとダウンロード
COPY go.mod ./
RUN go mod download

# ソースコードのコピー
COPY . .

# ビルド
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o mcp-meilisearch main.go

# 実行用の軽量イメージ
FROM alpine:latest

RUN apk --no-cache add ca-certificates

WORKDIR /root/

# ビルドされたバイナリのコピー
COPY --from=builder /app/mcp-meilisearch .
COPY config.yml .

EXPOSE 7701

# 実行
CMD ["./mcp-meilisearch", "server", "--config", "config.yml"]
