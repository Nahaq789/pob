# Battle Service 設計ドキュメント

## 概要

- バトル形式: 1v1 ターン制（擬似ターン制）
- 双方が同時に行動を選択し、素早さ・優先度順に片方ずつ処理する
- レベル固定: Lv.50 / メガシンカ・テラスタルなし

---

## アーキテクチャ

### WebSocket接続管理（SessionManager）

```go
type SessionManager struct {
    mu       sync.RWMutex
    sessions map[string][]*websocket.Conn // key: session_id
}
```

- `len(conns) == 2`：同一インスタンスに両プレイヤーがいる → 直接送信
- `len(conns) == 1`：相手が別インスタンスにいる → Redis Publish

---

## WebSocket

### エンドポイント

```
GET /ws/battle/:session_id
```

- ハンドシェイク時にAuthMiddlewareでJWT検証

### クライアント → サーバー

#### 技使用
```json
{
  "type": "move",
  "user_id": "uuid",
  "pokemon_id": 1,
  "move_id": 33
}
```

#### ポケモン交換
```json
{
  "type": "switch",
  "user_id": "uuid",
  "current_pokemon_id": 1,
  "next_pokemon_id": 2
}
```

#### サレンダー
```json
{
  "type": "surrender",
  "user_id": "uuid"
}
```

### サーバー → クライアント

#### 行動結果
```json
{
  "type": "action_result",
  "messages": [
    { "order": 1, "text": "ピカチュウのでんきショックだ!" },
    { "order": 2, "text": "こうかばつぐんだ!" }
  ],
  "players": {
    "player1": {
      "active_pokemon": {
        "box_pokemon_id": "uuid",
        "current_hp": 100,
        "pp": [15, 10, 5, 0],
        "status": "paralysis",
        "rank": {
          "attack": 0, "defense": -1, "sp_attack": 0,
          "sp_defense": 0, "speed": 0, "accuracy": 0, "evasion": 0
        },
        "is_alive": true
      },
      "party": [
        {
          "box_pokemon_id": "uuid",
          "current_hp": 100,
          "pp": [15, 10, 5, 0],
          "status": "",
          "is_alive": true
        }
      ]
    },
    "player2": {}
  },
  "is_finished": false
}
```

#### エラー通知
```json
{
  "type": "error",
  "message": "PPが残っていません"
}
```

---

## タイムアウト

- 制限時間: 45秒
- クライアント側: タイムアウトUIを表示し、時間切れでサレンダー送信
- サーバー側: `time.AfterFunc`で45秒タイマーを管理し、時間切れで未選択プレイヤーをサレンダー処理（二重管理）
- 両プレイヤーの選択が揃ったらタイマーをStop

---

## 接続切れ

- 再接続: なし
- 切断時: Redisのsession関連キーを削除し、相手に切断通知を送信
- 検知: pingpongで生存確認

---

## Redisステート設計

### キー構成

| キー | 内容 | 更新タイミング |
|------|------|------|
| `battle:{session_id}:static` | 固定データ | バトル開始時のみ |
| `battle:{session_id}:dynamic` | 動的データ | 毎ターン |

### static

```json
{
  "players": {
    "player1": {
      "user_id": "uuid",
      "party": [
        {
          "box_pokemon_id": "uuid",
          "slot": 1,
          "pokemon_id": 1,
          "name": "フシギダネ",
          "type1_id": 12,
          "type2_id": 0,
          "ability_id": 65,
          "nature": "adamant",
          "max_hp": 155,
          "attack": 123,
          "defense": 95,
          "sp_attack": 80,
          "sp_defense": 95,
          "speed": 110,
          "moves": [
            {
              "move_id": 33,
              "name": "たいあたり",
              "type_id": 1,
              "damage_class": "physical",
              "power": 40,
              "accuracy": 100,
              "pp": 15
            }
          ]
        }
      ]
    },
    "player2": {}
  }
}
```

### dynamic

```json
{
  "phase": "selecting",
  "acting_player": "",
  "pending_actions": {
    "player1": null,
    "player2": null
  },
  "players": {
    "player1": {
      "active_slot": 1,
      "party": [
        {
          "box_pokemon_id": "uuid",
          "current_hp": 155,
          "pp": [15, 10, 5, 0],
          "rank": {
            "attack": 0, "defense": 0, "sp_attack": 0,
            "sp_defense": 0, "speed": 0, "accuracy": 0, "evasion": 0
          },
          "status": "",
          "held_item_id": 0
        }
      ]
    },
    "player2": {}
  }
}
```

### phase一覧

| 値 | 説明 |
|------|------|
| `selecting` | 双方の行動選択待ち |
| `acting` | 行動処理中（acting_playerが行動中） |
| `turn_end` | 双方の行動終了後の処理中 |
| `finished` | バトル終了 |

### pending_actions

```go
type ActionType string

const (
    ActionTypeMove     ActionType = "move"
    ActionTypeSwitch   ActionType = "switch"
    ActionTypeSurrender ActionType = "surrender"
)

type PlayerActionDetail interface {
    GetType() ActionType
}

type PlayerAction struct {
    Detail PlayerActionDetail
}

type PendingActions struct {
    Actions map[string]*PlayerAction // key: "player1" | "player2"
}

type MoveAction struct {
    MoveId int
}
func (m *MoveAction) GetType() ActionType { return ActionTypeMove }

type SwitchAction struct {
    Slot int
}
func (s *SwitchAction) GetType() ActionType { return ActionTypeSwitch }

// Redis保存用
type PendingActionJSON struct {
    Type   ActionType `json:"type"`
    MoveId int        `json:"move_id,omitempty"`
    Slot   int        `json:"slot,omitempty"`
}
```

---

## バトルフェーズ

| フェーズ | タイミング | 例 |
|------|------|------|
| 1 | 場に出たとき | 威嚇・天候変化系特性 |
| 2 | 技選択前 | こだわり系持ち物・ポケモン交換選択 |
| 3 | ダメージ計算前 | ○○スキン系特性・ジュエル系持ち物・溜め技 |
| 4 | ダメージ計算中 | たいねつ・ふかしのこぶし・連続攻撃 |
| 5 | ダメージ計算後 | じきゅうりょく・追加効果・変化技 |
| 6 | 双方攻撃終了後 | 状態異常スリップ・天候スリップ・加速 |

---

## ハンドラーパイプライン

| Tier | 種別 |
|------|------|
| A | 汎用ハンドラー（多くの技・特性に共通） |
| B | 汎用ハンドラー（複数共有できる特殊処理） |
| C | 個別ハンドラー（固有処理） |

- 各タイミングにハンドラー配列を用意し優先度でソート
- レジストリ方式で管理

---

## ダメージ計算

- 乱数: 0.85〜1.00の16段階からランダム選択
- 実数値: `pkg/stats`で計算済み（nature補正済み）をstaticに保存
- ダメージ計算時: `実数値 × ランク補正テーブル`で算出
- タイプ相性テーブル: battle-serviceに静的データとして保持（18×18）

---

## 未決定事項

- [ ] バトル終了エンドポイントの設計
- [ ] Redis Pub/Subの詳細設計
- [ ] アトミック性の担保（SETNX / Lua script）の詳細
