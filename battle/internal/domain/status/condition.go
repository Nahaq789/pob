package status

import "context"

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

type OtherConditionHandler interface {
	Execute(ctx context.Context, s OtherStatus) error
}
