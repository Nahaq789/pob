#!/bin/bash

set -e

echo "=== dex-service 動作確認 ==="

cd "$(dirname "$0")/.."

if [ -f .env ]; then
    export $(grep -v '^#' .env | xargs)
else
    echo ".envファイルが見つかりません"
    exit 1
fi

go build -o /tmp/test-dex ./dex/cmd/test/main.go
/tmp/test-dex
rm /tmp/test-dex

echo "=== 完了 ==="