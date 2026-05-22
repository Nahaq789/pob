# Architecture

## Monorepo

`go.work` で全モジュールを管理。各サービスは独立した `go.mod` を持つ。

```
pob/
├── proto/          # .proto定義のみ（言語非依存）
├── pkg/            # サービス横断共有 (module: pob/pkg)
│   ├── logger/
│   ├── jwt/        # VerifyToken, ExtractBearerToken
│   ├── tracing/
│   ├── stats/      # 実数値計算（nature補正済み）
│   └── redis/
├── dex/            # module: pob/dex
├── user/           # module: pob/user
├── box/            # module: pob/box
├── battle/         # module: pob/battle
└── sync/           # module: pob/sync（PokeAPI同期バッチ）
```

## Services

| サービス | REST | gRPC | DB | Redis |
|----------|------|------|----|-------|
| dex | :8081 | :9091 | dex_db | ✅ |
| user | :8082 | — | user_db | ❌ |
| box | :8083 | :9093 | box_db | ❌ |
| battle | :8084 | — | なし | ✅ |

## gRPC Connections

```
box    → dex
battle → dex
battle → box
```

生成コードの配置:
- `dex/gen/`, `box/gen/` — サーバ用 + クライアント用
- `battle/internal/gen/` — クライアント用のみ

## Internal Structure

### dex / user / box（3層構造）

```
<service>/internal/
├── handler/        # HTTPハンドラ / gRPCサーバ実装
├── service/        # ユースケース
│   └── dto/        # リクエスト/レスポンス定義
├── repository/     # DB・外部アクセス
├── model/          # ドメインモデル
│   └── apperror/   # アプリケーションエラー定義
├── middleware/
└── shared/         # DBクライアント・トランザクション管理
```

### battle（DDD構成）

```
battle/internal/
├── presentation/websocket/
├── application/
├── domain/
│   ├── battle/
│   ├── damage/
│   └── move/
└── infrastructure/redis/
```

## Notes

- `.proto` はサービス実装時に都度追加（upfront設計不要）
- `pkg/stats` は `battle` と `box` 両方からインポート（dex の責務外）
