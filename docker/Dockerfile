# ビルドステージ
FROM golang:1.24-bullseye AS builder

# 作業ディレクトリの設定
WORKDIR /build

# ソースコードをコピー
COPY . .

# mcp-meilisearchをビルド
RUN make deps && make bin/mcp-meilisearch

# 実行ステージ
FROM getmeili/meilisearch:v1.13

# 必要なパッケージをインストール
RUN apk add --no-cache bash

# バイナリをビルドステージからコピー
COPY --from=builder /build/bin/mcp-meilisearch /bin/mcp-meilisearch

# バイナリに実行権限を付与
RUN chmod +x /bin/mcp-meilisearch

# 設定テンプレートとスタートスクリプトをコピー
COPY docker/config.yml.template /etc/mcp-meilisearch/config.yml.template
COPY docker/entrypoint.sh /bin/entrypoint.sh

# スクリプトに実行権限を付与
RUN chmod +x /bin/entrypoint.sh

# データディレクトリの作成
VOLUME ["/data.ms"]

# 作業ディレクトリを設定
WORKDIR /

# エントリポイントを設定
ENTRYPOINT ["/bin/entrypoint.sh"]
