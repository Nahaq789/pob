# POB (Pokémon Online Battle) — Agent Rules

1v1ターン制ポケモンバトルゲーム。ローカル運用前提のマイクロサービス構成。

## Stack

| レイヤー      | 技術             |
| ------------- | ---------------- |
| Frontend      | Next.js          |
| Backend       | Go (Gin)         |
| Inter-service | gRPC             |
| Front ↔ Back  | REST / WebSocket |
| Cache         | Redis            |
| Auth          | JWT (RS256)      |

## Commands

```bash
docker compose up                                      # 全サービス起動
docker compose --profile tools run --rm sync          # syncバッチ単体実行
docker compose build <service>                        # 特定サービスのみビルド
```

## Architecture

`go.work` で全モジュールを管理。各サービスは独立した `go.mod` を持つ。

```
pob/
├── proto/          # .proto定義のみ
├── pkg/            # サービス横断共有 (module: pob/pkg)
│   ├── logger/
│   ├── jwt/        # VerifyToken, ExtractBearerToken
│   ├── tracing/
│   ├── interceptor/hmac/  # HmacClientInterceptor / HmacServerInterceptor
│   ├── redis/
│   └── stats/      # 実数値計算（nature補正済み）
├── dex/            # gRPCのみ(:9091)
├── user/           # REST(:8082)
├── box/            # REST(:8083) + gRPC(:9093)
├── battle/         # REST+WebSocket(:8084) / DDD構成
└── sync/           # PokeAPIバッチ（Cobra CLI）
```

### Internal structure: dex / user / box

```
<service>/internal/
├── handler/        # HTTP or gRPC
├── service/
│   └── dto/
├── repository/
├── model/
│   └── apperror/
├── middleware/
└── shared/         # DBクライアント
```

### Internal structure: battle (DDD)

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

## Auth Design

- RS256、鍵ペアは1組のみ
- 秘密鍵: `user/pem/private.pem`（user-serviceのみ保持）
- 公開鍵: 各サービスの環境変数で配布
- ペイロード: `user_id`, `exp`, `iat` のみ
- アクセストークン: 15分、レスポンスボディ返却
- リフレッシュトークン: 7日、HttpOnly Cookie返却
- gRPCサービス間: HMAC共有シークレットをInterceptorで検証（`pkg/interceptor/hmac/`）
- `refresh_tokens` テーブルは1ユーザー1レコード（UPSERT）
- ログアウト: DELETE、`RowsAffected == 0` で `ErrAlreadyLoggedOut`

## Database

### user_db

- `users`: id(UUID PK), username(UNIQUE), password_hash
- `refresh_tokens`: id, user_id(UNIQUE FK), token_hash, expires_at

### dex_db

- `pokemon`, `types`, `abilities`, `pokemon_abilities`, `moves`, `pokemon_moves`, `items`
- PK は PokeAPI の INT ID
- `power` / `accuracy` は NULL許容（必中・固定ダメージ技）
- 性格（nature）補正は **Go側静的定義**（DBテーブルなし）

### box_db

- `boxes`, `box_pokemon`（IVs/EVs/技スロット1〜4 NULL許容）
- `parties`, `party_pokemon`（slot 1〜6, UNIQUE(party_id, slot)）
- ボックス上限（30匹）はアプリケーション層で制御
- PKの型: UUID（user/box系）

## Coding Rules

### Hashing

パスワード・リフレッシュトークンは SHA-256 → bcrypt の2段階ハッシュ。

```go
repository.Hash(plain string) (string, error)
repository.Compare(hashed, plain string) bool
```

### Logging

`log/slog` を使用。必ず Context を渡す。

```go
slog.InfoContext(ctx, "message", slog.String("key", val))
slog.ErrorContext(ctx, "message", slog.Any("error", err))
```

### Error Handling

- アプリケーションエラーは `model/apperror/` に定義
- ハンドラーでは `errors.Is` で判定してHTTPステータスを分岐

### Transaction

```go
txCtx := shared.TxWithContext(ctx, tx)
tx := shared.TxFromContext(ctx)  // repository層で取り出す
```

`TransactionManager.WithTransaction` 経由で使用。

### DB Client

- user / box / sync: `pgxpool`
- dex: `pgxpool` + GORM

### Middleware

リクエストには必ず `TraceMiddleware()` を適用（router登録時に `Use`）。

## Battle Design

### Redis State（2キー構成）

- `battle:{session_id}:static`: 固定データ（実数値/技/タイプ等）、バトル開始時のみ書き込み
- `battle:{session_id}:dynamic`: 動的データ（HP/PP/ランク/状態異常等）、毎ターン更新

### Damage Calculation

- フロント → battle-service: 技IDのみ送信
- battle-service がRedisから状態取得して計算完結
- 乱数: 0.85〜1.00の16段階
- タイプ相性テーブル: battle-service に静的データとして保持（18×18）

### Handler Pipeline

介入タイミング5種: 技使用前 / ダメージ計算前 / ダメージ計算中 / ダメージ計算後 / ターン終了時

| Tier | 種別                                     |
| ---- | ---------------------------------------- |
| A    | 汎用ハンドラー（多くの技・特性に共通）   |
| B    | 汎用ハンドラー（複数共有できる特殊処理） |
| C    | 個別ハンドラー（固有処理）               |

レジストリ方式で管理。全技実装は不要。拡張可能なアーキテクチャを示すことが目的。
