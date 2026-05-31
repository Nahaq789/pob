#!/bin/bash

set -e

echo "=== POB 環境構築 ==="

# RSA鍵ペア生成
echo "RSA鍵ペアを生成します..."
mkdir -p user/pem
openssl genrsa -out user/pem/private.pem 2048
openssl rsa -in user/pem/private.pem -pubout -out user/pem/public.pem
echo "RSA鍵ペア生成完了"

# HMACシークレット生成
echo "HMACシークレットを生成します..."
HMAC_SECRET=$(openssl rand -hex 32)
echo "HMACシークレット生成完了"

# docker-compose.override.yml生成
echo "docker-compose.override.ymlを生成します..."
cat > docker-compose.override.yml << EOF
services:
  dex:
    environment:
      - DB_DSN=postgres://pob:pob@dex-postgres:5432/dex_db?sslmode=disable
      - REDIS_ADDR=redis:6379
      - HMAC_SECRET=${HMAC_SECRET}

  user:
    environment:
      - DB_DSN=postgres://pob:pob@user-postgres:5432/user_db?sslmode=disable
      - PRIVATE_KEY_PATH=/pem/private.pem
      - PUBLIC_KEY_PATH=/pem/public.pem

  box:
    environment:
      - DB_DSN=postgres://pob:pob@box-postgres:5432/box_db?sslmode=disable
      - DEX_GRPC_ADDR=dex:9091
      - HMAC_SECRET=${HMAC_SECRET}

  battle:
    environment:
      - REDIS_ADDR=redis:6379
      - DEX_GRPC_ADDR=dex:9091
      - BOX_GRPC_ADDR=box:9093
      - HMAC_SECRET=${HMAC_SECRET}

  sync:
    environment:
      - DB_DSN=postgres://pob:pob@dex-postgres:5432/dex_db?sslmode=disable
      - POKEAPI_BASE_URL=https://pokeapi.co/api/v2

  user-postgres:
    environment:
      - POSTGRES_USER=pob
      - POSTGRES_PASSWORD=pob
      - POSTGRES_DB=user_db

  dex-postgres:
    environment:
      - POSTGRES_USER=pob
      - POSTGRES_PASSWORD=pob
      - POSTGRES_DB=dex_db

  box-postgres:
    environment:
      - POSTGRES_USER=pob
      - POSTGRES_PASSWORD=pob
      - POSTGRES_DB=box_db
EOF
echo "docker-compose.override.yml生成完了"

# .env生成部分を追加
cat > .env << EOF
HMAC_SECRET=${HMAC_SECRET}
EOF

echo "=== 環境構築完了 ==="