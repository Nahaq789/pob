#!/bin/bash

set -e

echo "=== DBデータ削除 ==="

cd "$(dirname "$0")/.."

echo "--- dex_db ---"
docker compose exec dex-postgres psql -U pob -d dex_db -c "TRUNCATE TABLE pokemon_moves, pokemon_abilities, moves, abilities, pokemon, types, items CASCADE;"

echo "--- user_db ---"
docker compose exec user-postgres psql -U pob -d user_db -c "TRUNCATE TABLE refresh_tokens, users CASCADE;"

echo "--- box_db ---"
docker compose exec box-postgres psql -U pob -d box_db -c "TRUNCATE TABLE party_pokemon, parties, box_pokemon, boxes CASCADE;"

echo "=== 削除完了 ==="
