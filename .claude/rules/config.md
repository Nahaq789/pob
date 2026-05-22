# Coding Rules

## Hashing

パスワード・リフレッシュトークンは SHA-256 → bcrypt の2段階ハッシュ。

```go
// user/internal/repository/hash.go
repository.Hash(plain string) (string, error)
repository.Compare(hashed, plain string) bool
```

## Logging

`log/slog` を使用。必ず Context を渡すこと。

```go
slog.InfoContext(ctx, "message", slog.String("key", val))
slog.ErrorContext(ctx, "message", slog.Any("error", err))
```

TraceID はミドルウェアで context に注入済み（`pkg/tracing`）。

## Error Handling

- アプリケーションエラーは `model/apperror/` に定義
- ハンドラーでは `errors.Is` で判定してHTTPステータスを分岐

```go
if errors.Is(err, apperror.ErrInvalidCredentials) {
    ctx.JSON(http.StatusUnauthorized, ...)
}
```

## Transaction

```go
// 開始
txCtx := shared.TxWithContext(ctx, tx)

// 取得（repository層でcontextから取り出す）
tx := shared.TxFromContext(ctx)
```

`TransactionManager.WithTransaction` 経由で使用する。

## DB Client

- user / box / sync: `pgxpool`（`user/internal/shared/db_client.go`）
- dex: `pgxpool` + `gorm`（`dex/internal/shared/db_client.go`）

## Middleware

リクエストには必ず `TraceMiddleware()` を適用すること（router登録時に `Use`）。
