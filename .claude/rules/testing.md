# Testing

## 基本方針

- 外部ライブラリ（testify等）は使用しない。標準 `testing` パッケージのみ使用。
- テーブル駆動テスト・サブテストを積極的に使う。
- テストは実装の観察可能な振る舞いを検証する（内部状態は原則触らない）。

---

## パッケージ構成

| 種別 | パッケージ名 | 用途 |
|------|-------------|------|
| 内部テスト | `package <pkg>` | 非公開フィールド・関数の検証が必要な場合 |
| 外部テスト | `package <pkg>_test` | 公開 API のみで完結する場合（原則こちら） |

```go
// 内部テスト例（rank パッケージ）
package rank

// 外部テスト例（rule パッケージ）
package rule_test
```

---

## テスト関数・サブテスト命名

### テスト関数名

```
Test<対象型またはコンストラクタ>_<メソッド名>
```

```go
func TestFlinch_Handle(t *testing.T) { ... }
func TestEncore_Resolve(t *testing.T) { ... }
func TestPreDamageHandle_Sleep(t *testing.T) { ... }
func TestNewBasicRank(t *testing.T) { ... }
```

### サブテスト名（`t.Run` の第1引数）

日本語で「**条件: 期待結果**」の形式で記述する。

```go
t.Run("残りターンあり: cleared=false", ...)
t.Run("封じられた技を選択: blocked=true", ...)
t.Run("アンコール技を選択: PhaseDamage に進む", ...)
t.Run("上限クランプ: delta=1", ...)
```

複数の観点を別サブテストに分割する（Phase とメッセージは分ける）。

```go
// good
t.Run("ひるみあり: PhaseEnd", ...)
t.Run("ひるみあり: ひるみメッセージあり", ...)

// bad
t.Run("ひるみあり", ...) // 何を検証するか不明
```

---

## サブテスト形式の選び方

### インライン `t.Run`（rule パッケージ準拠）

状態異常・技・特性など、**ケースごとにセットアップが異なる**場合に使う。

```go
func TestFlinch_Handle(t *testing.T) {
    t.Run("ひるみあり: PhaseEnd", func(t *testing.T) {
        fl := rule.NewFlinch()
        st := status.NewStatusWith(nil, []status.OtherStatus{fl})
        // ...
        if result.NextPhase != phase.PhaseEnd {
            t.Errorf("expected PhaseEnd, got %v", result.NextPhase)
        }
    })

    t.Run("ひるみあり: ひるみメッセージあり", func(t *testing.T) {
        // ...
    })
}
```

### テーブル駆動テスト（rank パッケージ準拠）

入力と期待値の組み合わせが多く、**セットアップが共通**の場合に使う。

```go
func TestUp(t *testing.T) {
    tests := []struct {
        name         string
        initialStage int
        up           int
        wantStage    int
        wantValue    float64
        wantMsg      string
    }{
        {"1段上昇", 0, 1, 1, 1.5, "が あがった"},
        {"2段上昇", 0, 2, 2, 2.0, "が ぐーんとあがった"},
        {"上限クランプ: delta=1", 5, 3, 6, 4.0, "が あがった"},
    }
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            // ...
        })
    }
}
```

---

## アサーション

### 基本パターン

標準の `if` ＋ `t.Error` / `t.Errorf` を使う。

```go
if result.NextPhase != phase.PhaseEnd {
    t.Errorf("expected PhaseEnd, got %v", result.NextPhase)
}
if msg == "" {
    t.Error("expected non-empty message")
}
if msg != "" {
    t.Errorf("expected empty message, got: %v", msg)
}
```

### セットアップ失敗時

テスト本体の実行が不可能な場合は `t.Fatal` / `t.Fatalf` で即終了する。

```go
ms, err := status.NewSleep(vo.NewCount(2))
if err != nil {
    t.Fatal(err)
}

got, err := r.Up(tt.up)
if err != nil {
    t.Fatalf("unexpected error: %v", err)
}
```

### 文字列部分一致

メッセージの完全一致が不要な場合は `strings.Contains` ベースのヘルパーを使う。

```go
func containsMsg(msgs []string, sub string) bool {
    for _, m := range msgs {
        if strings.Contains(m, sub) {
            return true
        }
    }
    return false
}

if !containsMsg(result.Messages, "ひるんで") {
    t.Errorf("expected flinch message, got %v", result.Messages)
}
```

---

## ランダム性を含む処理のテスト

`math/rand` を使う処理（まひ行動不能・こんらん自傷等）は以下の2パターンで検証する。

### パターン1: 両極端が発生することを確認

```go
t.Run("true と false（行動不能/行動可）の両方が返る", func(t *testing.T) {
    gotEnd, gotDamage := false, false
    for range 300 {
        // ... セットアップ ...
        result := h.Handle(ctx)
        switch result.NextPhase {
        case phase.PhaseEnd:
            gotEnd = true
        case phase.PhaseDamage:
            gotDamage = true
        }
        if gotEnd && gotDamage {
            break
        }
    }
    if !gotEnd {
        t.Error("300試行で行動不能（PhaseEnd）が一度も発生しなかった")
    }
    if !gotDamage {
        t.Error("300試行で行動成功（PhaseDamage）が一度も発生しなかった")
    }
})
```

### パターン2: 発生率が許容範囲内であることを確認

```go
t.Run("自傷率がおよそ 1/3 である", func(t *testing.T) {
    const trials = 3000
    hits := 0
    for range trials {
        if _, hit := c.CheckSelfHit("テスト"); hit {
            hits++
        }
    }
    rate := float64(hits) / trials
    if rate < 0.20 || rate > 0.47 {
        t.Errorf("unexpected self-hit rate: %.2f (expected ~0.33)", rate)
    }
})
```

試行回数の目安：

| 確率 | 試行回数 |
|------|---------|
| 1/3（こんらん自傷） | 300 |
| 1/8（まひ行動不能） | 300 |
| 確率分布の検証 | 3000 |

---

## ヘルパー関数

テストファイル内にローカルヘルパーを定義して、セットアップの重複を排除する。

```go
// ポケモン生成ヘルパー
func newTestPokemon(name string, speed int) *pokemon.Pokemon { ... }

// プレイヤー生成ヘルパー（選出まで完了させた状態で返す）
func newTestPlayer(id, name string, poke0, poke1, poke2 *pokemon.Pokemon) *player.Player { ... }
```

命名規則：`new<対象型>` プレフィックスを使い、テストファイル内の役割が分かる名前にする。

複数テストファイルが同じパッケージに存在する場合、ヘルパー名が衝突しないよう
テスト対象を示すプレフィックスを付ける（例: `newPreDamagePokemon`）。

---

## コメント

コメントはデフォルト不要。以下の場合のみ記載する。

- **クランプ計算の境界値**：なぜその入力でその出力になるか自明でない場合
  ```go
  // stage=5 + up=3 → 8 → clamp 6, delta=1
  ```
- **実装上の制約**：テスト対象の設計的な制約を説明する場合
  ```go
  // NewMainStatus(Freeze) は count=0 で生成されるため IsFrozen が即時 false を返し、解凍扱いとなる。
  ```

---

## 禁止事項

- `t.Skip` による省略（代わりにランダムテストの方針で対処する）
- `time.Sleep` の使用
- グローバル変数や `init()` によるテスト間の状態共有
- 外部ライブラリへの依存（`github.com/stretchr/testify` 等）
