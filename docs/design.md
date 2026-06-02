# POB（Pokémon Online Battle）設計ドキュメント

## 概要

| 項目 | 内容 |
|------|------|
| プロジェクト名 | POB（Pokémon Online Battle） |
| バトル形式 | 1対1 ターン制 |
| 対戦形式 | プレイヤー vs プレイヤー（異なるブラウザ） |
| 運用 | ローカル運用前提（デプロイなし） |
| レベル | 固定 Lv.50 |
| 非対応要素 | テラスタル・メガシンカなし |
| 外部API | PokeAPI 連携 |

---

## 技術スタック

| レイヤー | 技術 |
|----------|------|
| フロントエンド | Next.js |
| バックエンド | Go（Gin） |
| サービス間通信 | gRPC |
| フロント ↔ バック | REST / WebSocket |
| キャッシュ | Redis |
| 認証 | JWT（RS256） |

アーキテクチャ：**マイクロサービス構成**

---

## ディレクトリ構成（モノレポ）

```
pob/
├── web/                  # Next.js フロントエンド
├── dex/                  # dex-service
├── user/                 # user-service
├── box/                  # box-service
├── battle/               # battle-service（DDD構成）
├── proto/                # .protoファイル
├── sync/                 # PokeAPIデータ同期バッチ
├── pkg/                  # 共通パッケージ（logger・JWT・Redis・HMAC interceptor・stats）
└── docker-compose.yml
```

### 各サービスの内部構造

**dex / user / box**（3層構造）
```
<service>/
├── cmd/
├── internal/
│   ├── handler/
│   ├── service/
│   ├── repository/
│   └── model/
├── go.mod
└── go.sum
```

**battle**（DDD構成）
```
battle/
├── cmd/
├── internal/
│   ├── presentation/
│   │   └── websocket/
│   ├── application/
│   ├── domain/
│   │   ├── battle/
│   │   ├── damage/
│   │   └── move/
│   └── infrastructure/
│       └── redis/
├── go.mod
└── go.sum
```

---

## サービス構成

| サービス | 役割 | REST | gRPC | DB | Redis |
|----------|------|------|------|----|-------|
| dex-service | PokeAPIプロキシ・ポケモンデータ管理 | — | :9091 | dex_db | ✅ |
| user-service | 認証・JWT発行 | :8082 | — | user_db | ❌ |
| box-service | ボックス・編成管理 | :8083 | :9093 | box_db | ❌ |
| battle-service | バトル進行・WebSocket | :8084 | — | なし | ✅ |

## gRPC接続

```
box    → dex  (client→server)
battle → dex  (client→server)
battle → box  (client→server)
```

---

## 認証設計

- JWT鍵ペアは **1組**
- **秘密鍵**：user-service のみ保持（`user/pem/private.pem`、`.gitignore`除外）
- **公開鍵**：全サービスの環境変数に配置
- **サービス間（gRPC）**：共有シークレット（HMAC）を Interceptor で検証（`pkg/interceptor/hmac/`）
- **JWTペイロード**：`user_id`・`exp`・`iat` のみ
- **アクセストークン**：有効期限15分・レスポンスボディで返却
- **リフレッシュトークン**：有効期限7日・HttpOnly Cookieで返却・SHA-256→bcrypt二段階ハッシュ

---

## バトル設計

### Redisのバトル状態

| 種別 | 内容 |
|------|------|
| 固定 | ベース実数値（性格補正済み）・タイプ・技リスト |
| 動的 | 現在HP・各能力ランク（-6〜+6）・PP・状態異常・道具・特性 |

- タイプ相性テーブルは battle-service に**静的データ**として保持
- 実数値計算は `pkg/stats` で一元管理（battle・boxの両方からimport）

### ダメージ計算

- Next.js → battle-service へのリクエストは**技IDのみ**
- battle-service が Redis から状態取得して計算を完結させる
- 乱数: 0.85〜1.00 の 16段階からランダム選択
- ダメージ計算時：`ベース実数値 × ランク補正テーブル` で算出

### バリデーション（battle-service）

1. 技IDがポケモンの技リストに存在するか
2. PPが残っているか
3. 自分のターンか

### PP管理

- アイテムによる回復なし
- PP消費は battle-service が計算しレスポンスで返す（通常: -1、Pressure特性: -2）
- PP = 0 の技はフロントで非活性化

### 技・特性・道具のハンドラー設計

介入タイミング（5種）:
1. 技使用前
2. ダメージ計算前
3. ダメージ計算中
4. ダメージ計算後
5. ターン終了時

各タイミングにハンドラー配列を用意し、優先度でソート（タイプ変更系を先に処理）。

| Tier | 種別 |
|------|------|
| A | 汎用ハンドラー（多くの技・特性に共通） |
| B | 汎用ハンドラー（複数共有できる特殊処理） |
| C | 個別ハンドラー（固有処理） |

レジストリ方式で管理。

---

## 未決定事項

- [ ] 特性の詳細設計
