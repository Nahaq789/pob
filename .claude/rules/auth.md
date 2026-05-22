# Auth Design

## JWT

- アルゴリズム: RS256、鍵ペアは1組のみ
- **秘密鍵**: `user/pem/private.pem`（`.gitignore` 除外済み・user-service のみ保持）
- **公開鍵**: `user/pem/public.pem`（各サービスの環境変数で配布）
- **ペイロード**: `user_id`, `exp`, `iat` のみ
- **アクセストークン**: 有効期限 15分、レスポンスボディで返却
- **リフレッシュトークン**: 有効期限 7日、HttpOnly Cookie で返却

## Refresh Token

- `refresh_tokens` テーブルは1ユーザー1レコード
- ログイン・リフレッシュ時は UPSERT でローテーション
- ログアウト時は DELETE、`RowsAffected == 0` で `ErrAlreadyLoggedOut` を返す

## Inter-service Auth

- gRPC サービス間: 共有シークレット（HMAC）を Interceptor で検証

## Endpoints

| Method | Path | 説明 |
|--------|------|------|
| POST | /user/register | ユーザー登録 |
| POST | /auth/login | ログイン |
| POST | /auth/refresh | トークンリフレッシュ |
| POST | /auth/logout | ログアウト |
