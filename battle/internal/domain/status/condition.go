package status

type Condition string

const (
	Poison    Condition = "poison"
	BadPoison Condition = "bad_poison"
	Paralysis Condition = "paralysis"
	Sleep     Condition = "sleep"
	Burn      Condition = "burn"
	Freeze    Condition = "freeze"
	None      Condition = "none"
)

type OtherCondition string

const (
	Confusion OtherCondition = "confusion"
)

// 技が使用されたターンが終了することで消滅する
type ClearedOnMoveTurnEnd string

// 次の行動時に消滅する
type ClearedOnNextAction string

// ターン経過で消滅する
type ClearedOverTurns string

// その他の条件で消滅する
type ClearedOnOtherCondition string

// 交代しない限り永続
type PersistUntilSwitch string

// あばれるけい
type Rampage string

// 眠り状態から回復したとき
type ClearedOnWakeUp string

// でんき技で行動したとき（Gen9仕様を採用。Gen1には該当技なし）
type ClearedOnElectricMove string
