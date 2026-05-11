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
├── pkg/                  # 共通パッケージ（logger・JWT）
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
├── shared/               # DB・Redisクライアント
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
├── shared/
├── go.mod
└── go.sum
```

---

## サービス構成

| サービス | 役割 | DB | Redis |
|----------|------|----|-------|
| dex-service | PokeAPIプロキシ・ポケモンデータ管理 | Pokemon DB | ✅ |
| user-service | 認証・JWT発行 | User DB | ❌ |
| box-service | ボックス・編成管理 | Box DB | ❌ |
| battle-service | バトル進行・WebSocket | なし | ✅ |

---

## 認証設計

- JWT鍵ペアは **1組**
- **秘密鍵**：user-service のみ保持
- **公開鍵**：全サービスの環境変数に配置
- **サービス間（gRPC）**：共有シークレット（HMAC）を Interceptor で検証
- **JWTペイロード**：`user_id`・`exp`・`iat` のみ

---

## バトル設計

### Redisのバトル状態

| 種別 | 内容 |
|------|------|
| 固定データ | ベース実数値（性格補正済み）・タイプ・技リスト |
| 動的データ | 現在HP・各能力ランク（-6〜+6）・PP・状態異常・道具・特性 |

- タイプ相性テーブルは battle-service に**静的データ**として保持

### ダメージ計算

- 乱数：battle-service で **0.85〜1.00 の16段階**からランダム選択
- 実数値：dex-service が**性格補正済み**で計算・返却
- ダメージ計算時：`ベース実数値 × ランク補正テーブル` で算出
- Next.js → battle-service へのリクエストは**技IDのみ**
- battle-service は Redis から状態取得して計算完結

### バリデーション（battle-service）

1. 技IDがポケモンの技リストに存在するか
2. PPが残っているか
3. 自分のターンか

### PP管理

- アイテムによる回復なし
- PP消費は battle-service が計算しレスポンスで返す
  - 通常：`-1`
  - 相手特性「プレッシャー」時：`-2`
- PP = 0 の技はフロントで**非活性化**

### 技・特性・道具のハンドラー設計

**介入タイミング（5種）**

1. 技使用前
2. ダメージ計算前
3. ダメージ計算中
4. ダメージ計算後
5. ターン終了時

**処理方式**

- 各タイミングにハンドラー配列を用意
- 優先度でソート（タイプ変更系を先に処理）
- ダメージ計算用構造体をパイプライン的に処理

**ハンドラー種別**

| Tier | 種別 | 説明 |
|------|------|------|
| A | 汎用ハンドラー | 多くの技・特性に共通する処理 |
| B | 汎用ハンドラー | やや特殊だが複数で共有できる処理 |
| C | 個別ハンドラー | 固有の処理を持つ技・特性・道具 |

レジストリ方式で管理。

---

## 未決定事項

- [ ] APIエンドポイント設計（REST / gRPC）
- [ ] DBスキーマ設計
- [ ] `.proto` ファイル設計
- [ ] 特性の詳細設計
