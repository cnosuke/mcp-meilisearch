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

# supervisorのインストール (Alpine用)
RUN apk add --no-cache supervisor bash

# バイナリをビルドステージからコピー
COPY --from=builder /build/bin/mcp-meilisearch /bin/mcp-meilisearch

# バイナリに実行権限を付与
RUN chmod +x /bin/mcp-meilisearch && \
    ls -la /bin/mcp-meilisearch

# 設定テンプレートとスタートスクリプトをコピー
COPY docker/config.yml.template /etc/mcp-meilisearch/config.yml.template
COPY docker/start.sh /bin/start.sh

# スタートスクリプトに実行権限を付与
RUN chmod +x /bin/start.sh

# supervisorの設定ディレクトリを作成
RUN mkdir -p /etc/supervisor.d /var/log/supervisor

# supervisorの設定
COPY docker/supervisord.conf /etc/supervisord.conf

# データディレクトリの作成
VOLUME ["/data.ms"]

# 作業ディレクトリを設定
WORKDIR /

# スタートスクリプトを起動
CMD ["/bin/start.sh"]
