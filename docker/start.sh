#!/bin/sh
set -e

# APIキーが設定されていることを確認
if [ -z "${MEILI_MASTER_KEY}" ]; then
  echo "Error: MEILI_MASTER_KEY environment variable is required"
  echo "Please run with: -e MEILI_MASTER_KEY=your_secret_key"
  exit 1
fi

# APIキーの長さを確認 (最低16バイト)
if [ ${#MEILI_MASTER_KEY} -lt 16 ]; then
  echo "Error: MEILI_MASTER_KEY must be at least 16 characters long"
  echo "Current key length: ${#MEILI_MASTER_KEY} characters"
  echo "Please provide a longer master key"
  exit 1
fi

# ログディレクトリを作成
mkdir -p /var/log/supervisor

# APIキーを環境変数から設定
API_KEY=${MEILI_MASTER_KEY}

# mcp-meilisearch設定ファイルにAPIキーを設定
mkdir -p /etc/mcp-meilisearch
sed "s/{{API_KEY}}/$API_KEY/g" /etc/mcp-meilisearch/config.yml.template > /etc/mcp-meilisearch/config.yml

# デバッグ情報
echo "Current directory: $(pwd)"
echo "Meilisearch binary location: $(which meilisearch || echo 'Not found')"
echo "MCP-Meilisearch binary location: $(which mcp-meilisearch || echo 'Not found')"
ls -la /bin/mcp-meilisearch

# supervisordを起動
exec /usr/bin/supervisord -c /etc/supervisord.conf
