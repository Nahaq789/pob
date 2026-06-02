#!/bin/bash

set -e

echo "=== POB 環境構築 ==="

# RSA鍵ペア生成
echo "RSA鍵ペアを生成します..."
mkdir -p user/pem
openssl genrsa -out user/pem/private.pem 2048
openssl rsa -in user/pem/private.pem -pubout -out user/pem/public.pem
echo "RSA鍵ペア生成完了"

# 公開鍵をbox/pemにコピー
echo "公開鍵をbox/pemにコピーします..."
mkdir -p box/pem
cp user/pem/public.pem box/pem/public.pem
echo "公開鍵コピー完了"

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
      - PUBLIC_KEY_PATH=/pem/public.pem

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

# 各サービスのローカル開発用 .env 生成
echo "各サービスの .env を生成します..."

cat > dex/.env << EOF
DB_DSN=postgres://pob:pob@localhost:5433/dex_db?sslmode=disable
REDIS_ADDR=localhost:6379
HMAC_SECRET=${HMAC_SECRET}
EOF

cat > user/.env << EOF
DB_DSN=postgres://pob:pob@localhost:5432/user_db?sslmode=disable
PRIVATE_KEY_PATH=./pem/private.pem
PUBLIC_KEY_PATH=./pem/public.pem
EOF

cat > box/.env << EOF
DB_DSN=postgres://pob:pob@localhost:5434/box_db?sslmode=disable
DEX_GRPC_ADDR=localhost:9091
HMAC_SECRET=${HMAC_SECRET}
PUBLIC_KEY_PATH=./pem/public.pem
EOF

cat > battle/.env << EOF
REDIS_ADDR=localhost:6379
DEX_GRPC_ADDR=localhost:9091
BOX_GRPC_ADDR=localhost:9093
HMAC_SECRET=${HMAC_SECRET}
EOF

cat > sync/.env << EOF
DB_DSN=postgres://pob:pob@localhost:5433/dex_db?sslmode=disable
POKEAPI_BASE_URL=https://pokeapi.co/api/v2
EOF

echo ".env 生成完了"

echo "=== 環境構築完了 ==="