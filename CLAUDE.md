# POB (Pokémon Online Battle)

1v1ターン制ポケモンバトルゲーム。ローカル運用前提のマイクロサービス構成。

- バトル形式: 1v1 ターン制 PvP（異なるブラウザ）
- レベル固定: Lv.50 / メガシンカ・テラスタルなし
- 外部API: PokeAPI

## Stack

| レイヤー | 技術 |
|----------|------|
| Frontend | Next.js |
| Backend | Go (Gin) |
| Inter-service | gRPC |
| Front ↔ Back | REST / WebSocket |
| Cache | Redis |
| Auth | JWT (RS256) |

## Commands

```bash
docker compose up                                      # 全サービス起動
docker compose --profile tools run --rm sync          # syncバッチ単体実行
docker compose build <service>                        # 特定サービスのみビルド
```

## Rules

- @.claude/rules/architecture.md
- @.claude/rules/auth.md
- @.claude/rules/database.md
- @.claude/rules/battle.md
- @.claude/rules/coding.md
