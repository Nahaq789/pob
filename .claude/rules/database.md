# Database

## user_db

| テーブル | 主要カラム |
|----------|-----------|
| `users` | id (UUID PK), username (UNIQUE), password_hash |
| `refresh_tokens` | id, user_id (UNIQUE FK), token_hash, expires_at |

## dex_db

| テーブル | 備考 |
|----------|------|
| `pokemon` | PK は PokeAPI の INT ID |
| `types` | 同上 |
| `abilities` | 同上 |
| `pokemon_abilities` | slot, is_hidden |
| `moves` | power/accuracy は NULL許容（必中・固定ダメ技） |
| `pokemon_moves` | 複合PK |

- 性格（nature）補正は**静的 Go 定義**で処理（DBテーブルなし）

## box_db

| テーブル | 備考 |
|----------|------|
| `boxes` | user_id, name |
| `box_pokemon` | IVs/EVs/技スロット1〜4（NULL許容） |
| `parties` | user_id, name |
| `party_pokemon` | slot 1〜6、UNIQUE(party_id, slot) |

- ボックス上限（30匹）はアプリケーション層で制御

## Notes

- PKの型: UUID（user/box系）、INT（dex系 = PokeAPI IDを流用）
- 各サービスのDBは**他サービスから直接参照しない**（必ずgRPC経由）
