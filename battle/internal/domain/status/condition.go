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
