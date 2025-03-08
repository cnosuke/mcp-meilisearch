#!/bin/sh
set -e

# ログディレクトリを作成
mkdir -p /var/log

# すべてのログメッセージを/var/log/container.logにリダイレクト
exec 3>&1  # 標準出力を保持
exec > /var/log/container.log 2>&1

# ここからの標準出力はすべてファイルにリダイレクト
echo "Starting container initialization at $(date)"

# APIキーが設定されていることを確認
if [ -z "${MEILI_MASTER_KEY}" ]; then
  echo "Error: MEILI_MASTER_KEY environment variable is required"
  exit 1
fi

# APIキーの長さを確認 (最低16バイト)
if [ ${#MEILI_MASTER_KEY} -lt 16 ]; then
  echo "Error: MEILI_MASTER_KEY must be at least 16 characters long"
  echo "Current key length: ${#MEILI_MASTER_KEY} characters"
  exit 1
fi

# APIキーを環境変数から設定
API_KEY=${MEILI_MASTER_KEY}

# mcp-meilisearch設定ファイルにAPIキーを設定
mkdir -p /etc/mcp-meilisearch
sed "s/{{API_KEY}}/$API_KEY/g" /etc/mcp-meilisearch/config.yml.template > /etc/mcp-meilisearch/config.yml

# Meilisearchをバックグラウンドで起動
echo "Starting Meilisearch in background..."
/bin/meilisearch --db-path /data.ms --env production --http-addr 0.0.0.0:7700 --master-key "${MEILI_MASTER_KEY}" > /var/log/meilisearch.log 2> /var/log/meilisearch-error.log &
MEILI_PID=$!

# Meilisearchが起動するまで少し待機
echo "Waiting for Meilisearch to start..."
sleep 2

# Meilisearchプロセスが正常に起動しているか確認
if ! kill -0 $MEILI_PID 2>/dev/null; then
  echo "Error: Meilisearch failed to start. Check logs at /var/log/meilisearch-error.log"
  cat /var/log/meilisearch-error.log
  exit 1
fi

echo "Meilisearch started with PID $MEILI_PID"
echo "Starting MCP-Meilisearch with --no-logs option..."

# 標準出力を元に戻す（MCPプロトコル通信用）
exec 1>&3 3>&-

# MCP-Meilisearchをフォアグラウンドで実行（--no-logsオプション付き）
exec /bin/mcp-meilisearch server --config /etc/mcp-meilisearch/config.yml --no-logs --log /var/log/mcp-meilisearch.log
