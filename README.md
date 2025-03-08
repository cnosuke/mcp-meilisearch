# Meilisearch MCP Server (Go)

Model Context Protocol (MCP) サーバーとして動作する、Go言語で実装されたMeilisearch検索エンジンのラッパーです。

## 概要

このプロジェクトは、Meilisearch検索エンジンをMCP (Model Context Protocol) サーバーとして使用できるようにするためのGo言語アプリケーションです。MCPサーバーを介して、Claudeなどの大規模言語モデル（LLM）がMeilisearchの強力な検索機能にアクセスできるようになります。

## 特徴

- Go言語で実装された軽量かつ高速なMCPサーバー
- Meilisearch APIとのシームレスな統合
- MCPプロトコルに準拠したインターフェース
- 検索、インデックス管理、ドキュメント操作などの主要機能をサポート
- Docker対応

## インストール

### 前提条件

- Go 1.24以上
- 実行中のMeilisearchインスタンス

### ビルド方法

```bash
# リポジトリのクローン
git clone https://github.com/cnosuke/mcp-meilisearch.git
cd mcp-meilisearch

# 依存関係のインストール
make deps

# ビルド
make bin/mcp-meilisearch
```

### Dockerを使用する場合

```bash
# Dockerイメージのビルド
make docker-build

# 環境変数ファイルを作成（サンプルをコピー）
cp .env.example .env
# 必要に応じて.envファイルを編集

# Dockerでの実行
make docker-run
```

## 設定

`config.yml`ファイルを使用して、Meilisearchサーバーへの接続設定を行います：

```yaml
meilisearch:
  host: http://localhost:7700  # Meilisearchサーバーのアドレス
  api_key: ""                  # 必要に応じてAPIキーを設定
```

環境変数でも設定可能です：
- `MEILISEARCH_HOST`: Meilisearchサーバーのアドレス
- `MEILISEARCH_API_KEY`: Meilisearch APIキー

## 使用方法

### サーバーの起動

```bash
./bin/mcp-meilisearch server --config config.yml
```

オプション：
- `--no-logs`: ログ出力を最小化（エラーのみ表示）
- `--log <ファイルパス>`: 指定したファイルにログを出力

### MCPクライアントからの接続

このMCPサーバーは、Claude AI (Claude Desktop)などのMCP対応クライアントと連携できます。

Claude Desktopの場合の設定例：
```json
{
  "mcpServers": {
    "meilisearch": {
      "command": "/path/to/mcp-meilisearch",
      "args": ["server", "--config", "/path/to/config.yml"]
    }
  }
}
```

## 利用可能なツール

このMCPサーバーは以下のツールを提供します：

### サーバー管理

- `health_check`: Meilisearchサーバーの状態を確認します。
  - パラメータ: なし
  - 戻り値: サーバーのヘルスステータス情報

### インデックス管理

- `list_indexes`: 全てのインデックスのリストを取得します。
  - パラメータ: なし
  - 戻り値: インデックスの配列

- `create_index`: 新規インデックスを作成します。
  - パラメータ:
    - `uid`: インデックスの一意識別子（必須）
    - `primary_key`: ドキュメントのプライマリキー（オプション）
  - 戻り値: 作成タスクの情報

### ドキュメント操作

- `get_documents`: インデックスからドキュメントを取得します。
  - パラメータ:
    - `index_uid`: インデックスのUID（必須）
    - `limit`: 取得するドキュメントの最大数（オプション）
    - `offset`: スキップするドキュメント数（オプション）
    - `fields`: 取得するフィールドの配列（オプション）
  - 戻り値: ドキュメントの配列

- `add_documents`: インデックスにドキュメントを追加します。
  - パラメータ:
    - `index_uid`: インデックスのUID（必須）
    - `documents`: 追加するドキュメントの配列（必須）
    - `primary_key`: プライマリキー（オプション）
  - 戻り値: 追加タスクの情報

### 検索

- `search`: インデックス内のドキュメントを検索します。
  - パラメータ:
    - `index_uid`: 検索対象インデックスのUID（必須）
    - `query`: 検索クエリ（必須）
    - `limit`: 返される結果の最大数（オプション）
    - `offset`: スキップする結果数（オプション）
    - `filter`: フィルタ式（オプション）
    - `sort`: ソート基準の配列（オプション）
  - 戻り値: 検索結果

## 使用例

### インデックスの作成

```json
{
  "name": "create_index",
  "arguments": {
    "uid": "movies",
    "primary_key": "id"
  }
}
```

### ドキュメントの追加

```json
{
  "name": "add_documents",
  "arguments": {
    "index_uid": "movies",
    "documents": [
      {
        "id": 1,
        "title": "Carol",
        "genres": ["Romance", "Drama"]
      },
      {
        "id": 2,
        "title": "Wonder Woman",
        "genres": ["Action", "Adventure"]
      }
    ]
  }
}
```

### 検索の実行

```json
{
  "name": "search",
  "arguments": {
    "index_uid": "movies",
    "query": "wonder",
    "limit": 5
  }
}
```

## プロジェクト構造

```
mcp-meilisearch/
├── bin/                    # ビルド済みバイナリの出力先
├── config/                 # 設定関連のパッケージ
│   └── config.go           # 設定ロード機能
├── logger/                 # ロギング関連のパッケージ
│   └── logger.go           # ロガー初期化
├── server/                 # サーバー実装
│   ├── meilisearch.go      # Meilisearchクライアント管理
│   ├── server.go           # MCPサーバーのメイン実装
│   └── tools/              # ツール実装
│       ├── tools.go        # ツール登録
│       ├── health_check.go # ヘルスチェックツール
│       ├── list_indexes.go # インデックス一覧ツール
│       └── ...             # その他のツール
├── config.yml              # 設定ファイル
├── Dockerfile              # Dockerビルド設定
├── go.mod                  # Goモジュール定義
└── main.go                 # エントリーポイント
```

## 拡張方法

このプロジェクトは拡張しやすい構造になっています。新しいツールを追加するには：

1. `server/tools/`ディレクトリに新しいツール実装ファイルを作成します
2. `server/tools/tools.go`ファイルの`RegisterAllTools`関数内で新しいツールを登録します

## ライセンス

MIT

## 関連プロジェクト

- [Meilisearch](https://github.com/meilisearch/meilisearch) - 元となる検索エンジン
- [meilisearch-mcp](https://github.com/meilisearch/meilisearch-mcp) - Python版のMeilisearch MCP
- [meilisearch-go](https://github.com/meilisearch/meilisearch-go) - Go言語用Meilisearchクライアント
