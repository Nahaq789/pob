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

// 技が使用されたターンが終了することで消滅する
type ClearedOnMoveTurnEnd OtherCondition

const (
	Flinch ClearedOnMoveTurnEnd = "flinch"
)

// 次の行動時に消滅する
type ClearedOnNextAction OtherCondition

// ターン経過で消滅する
type ClearedOverTurns OtherCondition

const (
	MoveDisabled ClearedOverTurns = "move_disabled"
)

// その他の条件で消滅する
type ClearedOnOtherCondition OtherCondition

const (
	Confusion ClearedOnOtherCondition = "confusion"
)

// 交代しない限り永続
type PersistUntilSwitch OtherCondition

// あばれるけい
type Rampage OtherCondition

// 眠り状態から回復したとき
type ClearedOnWakeUp OtherCondition

// でんき技で行動したとき（Gen9仕様を採用。Gen1には該当技なし）
type ClearedOnElectricMove OtherCondition
