#!/bin/bash
set -e

echo "=== reset-db: start ==="

# 1. dex-postgres のテーブルを全削除して再作成
echo "[1/3] Recreating dex schema..."
docker compose exec -T dex-postgres psql -U pob -d dex_db << 'EOF'
DROP TABLE IF EXISTS pokemon_moves CASCADE;
DROP TABLE IF EXISTS pokemon_abilities CASCADE;
DROP TABLE IF EXISTS moves CASCADE;
DROP TABLE IF EXISTS abilities CASCADE;
DROP TABLE IF EXISTS pokemon CASCADE;
DROP TABLE IF EXISTS types CASCADE;
DROP TABLE IF EXISTS items CASCADE;
EOF

docker compose exec -T dex-postgres psql -U pob -d dex_db < scripts/dex-schema.sql

# 2. sync: CSV取得
echo "[2/3] Fetching CSV from PokeAPI..."
cd sync
go run ./... csv --target types
go run ./... csv --target items
go run ./... csv --target moves          --gen 1
go run ./... csv --target abilities      --gen 1
go run ./... csv --target pokemon        --gen 1
go run ./... csv --target pokemon_abilities --gen 1
go run ./... csv --target pokemon_moves  --gen 1

# 3. sync: DBインサート
echo "[3/3] Inserting data into DB..."
go run ./... sync --target types
go run ./... sync --target items
go run ./... sync --target moves          --gen 1
go run ./... sync --target abilities      --gen 1
go run ./... sync --target pokemon        --gen 1
go run ./... sync --target pokemon_abilities --gen 1
go run ./... sync --target pokemon_moves  --gen 1
cd ..

echo "=== reset-db: done ==="
