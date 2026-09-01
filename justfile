set dotenv-load := true

# ── Docker ────────────────────────────────────────────────────────────────────

# 全サービス起動
up:
    docker compose up

# 全サービス停止・コンテナ削除
down:
    docker compose down

# 特定サービスのみビルド（例: just build battle）
build service:
    docker compose build {{ service }}

# syncバッチ単体実行
sync:
    docker compose --profile tools run --rm sync

# ── Run（ローカル開発） ────────────────────────────────────────────────────────

@run-user:
    cd ./user && go run ./...

@run-dex:
    cd ./dex && go run ./...

@run-box:
    cd ./box && go run ./...

@run-battle:
    cd ./battle && go run ./...

# ── Test ──────────────────────────────────────────────────────────────────────

# 全モジュールのテスト
test:
    go test ./pkg/...
    cd ./user    && go test ./...
    cd ./dex     && go test ./...
    cd ./box     && go test ./...
    cd ./battle  && go test ./...

test-pkg:
    go test ./pkg/...

test-user:
    cd ./user && go test ./...

test-dex:
    cd ./dex && go test ./...

test-box:
    cd ./box && go test ./...

test-battle:
    cd ./battle && go test -cover ./...

# ── DB ────────────────────────────────────────────────────────────────────────

# dex_db を初期化して sync バッチ実行
reset-dex:
    bash ./scripts/reset-dex-db.sh

# 全テーブルをトランケート（データ削除、スキーマは保持）
truncate:
    bash ./scripts/truncate.sh

# ── Setup ────────────────────────────────────────────────────────────────────

# RSA 鍵ペア生成 + .env 作成
setup:
    bash ./scripts/setup.sh
