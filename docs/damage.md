# ダメージフェーズ 設計ドキュメント

## 概要

ダメージフェーズは以下2つのコンポーネントで構成される。

| コンポーネント | 場所 | 責務 |
|---|---|---|
| `DamagePhaseHandler` | `battle/internal/domain/phase/damage.go` | Battle状態を読み、計算に必要な情報を組み立てて進行を制御する |
| `damage.Calculate` | `battle/internal/domain/damage/damage.go` | 入力値を受け取り、ダメージ量（int）を返すだけ。副作用なし |

---

## damage.go — ダメージ計算

### 入力（`DamageInput` 構造体）

| フィールド | 型 | 説明 |
|---|---|---|
| `Power` | int | 技の威力 |
| `Attack` | int | 攻撃実数値（ランク補正・急所補正済み） |
| `Defense` | int | 防御実数値（ランク補正・急所補正済み） |
| `IsCrit` | bool | 急所フラグ |
| `Random` | int | 乱数（85〜100） |
| `IsStab` | bool | タイプ一致フラグ |
| `TypeEff` | float64 | タイプ相性係数（0 / 0.25 / 0.5 / 1.0 / 2.0 / 4.0 等） |
| `IsBurned` | bool | やけど状態フラグ |
| `Weather` | float64 | 天気補正係数（1.0 / 0.5 / 2.0） |
| `Wall` | float64 | 壁補正係数（1.0 / 0.5） |
| `Other` | float64 | 特性・アイテム補正の合成値 |

### 計算式（第五世代以降）

```
base = floor(floor(floor(Level×2/5+2) × Power × Attack / Defense) / 50) + 2
```

- Level は 50 固定（`(2×50/5+2) = 22`）
- 各ステップで端数切り捨て

### 修飾子 M の適用順

base に対して以下の補正を順番に乗算する。

| # | 補正 | 係数 | 備考 |
|---|------|------|------|
| 1 | 範囲補正 | 1.0 | 1v1 固定のため常に 1.0 |
| 2 | おやこあい補正 | 1.0 | メガシンカ未実装のため常に 1.0 |
| 3 | 天気補正 | `Weather` | 雨 + 水技 → 1.5倍、晴れ + 炎技 → 1.5倍 等 |
| 4 | 急所補正 | 1.5倍 or 1.0 | `IsCrit == true` のとき 1.5倍 |
| 5 | 乱数補正 | `Random / 100` | 85〜100 の 16段階 |
| 6 | STAB | 1.5倍 or 1.0 | `IsStab == true` のとき 1.5倍 |
| 7 | タイプ相性 | `TypeEff` | 0〜4倍 |
| 8 | やけど補正 | 0.5倍 or 1.0 | `IsBurned == true` かつ Physical のとき 0.5倍 |
| 9 | 壁補正 | `Wall` | リフレクター / ひかりのかべ |
| 10 | その他補正 | `Other` | 特性・アイテムの合成値 |

- 最終ダメージが 0 の場合は 1 にする

---

## DamagePhaseHandler — フェーズの進行制御

### DamageContext

```go
type DamageContext struct {
    Battle     *battle.Battle
    ActorId    string
    MoveId     int        // 通常時のみ有効。自傷時は無視
    Type       ptype.Type // 自傷時は ptype.None
    Power      int
    Category   move.DamageClass // Physical / Special
    MustHit    bool             // true なら命中判定スキップ
    CanCrit    bool             // false なら急所判定スキップ
    TargetSelf bool             // true なら Actor 自身が Target
}
```

### 処理フロー

#### 1. 命中判定（`MustHit == false` の場合のみ）
- Actor の命中ランク・Target の回避ランクから `AccuracyRank` を生成
- 技の `accuracy` と照合してロール
- **外れ → `NextPhase: PhaseEnd`**
- 外れた後に発動する挙動（例: じたんだ）は技ハンドラー側の責務

#### 2. 急所判定（`CanCrit == true` の場合のみ）
- Actor の急所ランクから確率を取得してロール
- `isCrit` フラグを以降の処理で保持

#### 3. A・D の組み立て
- `Category == Physical` → Attack / Defense を使用
- `Category == Special` → SpAttack / SpDefense を使用
- **急所時のランク無視ルール:** 攻撃側に不利なランク（マイナス）は無視、防御側に有利なランク（プラス）は無視

#### 4. 各係数の算出
| 係数 | 算出方法 |
|------|---------|
| `IsStab` | Actor のタイプ配列と `DamageContext.Type` を比較 |
| `TypeEff` | `ptype` のタイプ相性テーブルを参照（Target のタイプ × 技のタイプ） |
| `IsBurned` | Actor のメイン状態異常を確認 |
| `Weather` | フィールドの天気状態から係数を決定 |
| `Wall` | フィールドの壁状態から係数を決定 |
| `Other` | `DamageModifierHandler` 群を呼び出して乗算 |

#### 5. 乱数生成
- 85〜100 の 16段階からランダム選択

#### 6. `damage.Calculate(input)` を呼び出し

#### 7. HP 反映
- Target に対してダメージを適用

#### 8. 気絶チェック
- 気絶 → `NextPhase: PhaseEnd`
- 生存 → `NextPhase: PhasePostDamage`

---

## 特性・アイテムのダメージ補正ハンドラー

特性・アイテムが与える乗算補正は専用インターフェースで定義し、`DamagePhaseHandler` が収集して `Other` に合成する。

```go
type DamageModifierHandler interface {
    Modifier(ctx DamageContext) float64
}
```

- 各ハンドラーが係数（float64）を返す
- `Other = mod1 × mod2 × ...` として `DamageInput` に渡す

---

## PreDamagePhase との連携

### タイプ変更系特性（リベロ・へんげんじざい 等）

実際のゲームでも「ダメージ計算の前にタイプが確定する」ため、`PreDamagePhase` で `DamageContext.Type` を書き換える。

**PreDamagePhase の処理順（確定版）:**

```
1. ねむり / こおり 判定       → 行動不能なら PhaseEnd
2. ひるみ 判定               → 行動不能なら PhaseEnd
3. アンコール 判定            → 行動不能なら PhaseEnd
4. かなしばり 判定            → 行動不能なら PhaseEnd
5. こんらん 判定（自傷チェック）→ 自傷ダメージ発生なら PhaseEnd
──────────────────────────────────────────────────
6. タイプ変更系特性ハンドラー   ← DamageContext.Type を書き換える ※未実装
──────────────────────────────────────────────────
7. まひ 判定                 → 行動不能なら PhaseEnd
8. PhaseDamage へ
```

> **Note:** タイプ変更系特性（リベロ・へんげんじざい 等）は比較的新しい特性のため、現時点では未実装。構造上は 5 と 7 の間にハンドラーを差し込む形で拡張可能。

---

## スコープ外（今回実装しない）

| 項目 | 理由 |
|------|------|
| おやこあい補正 | メガシンカ未実装（メガガルーラ専用） |
| 範囲補正 | 1v1 固定 |
| ダイマックス関連補正 | 未実装 |
| タイプ変更系特性 | 比較的新しい特性のため後回し |
