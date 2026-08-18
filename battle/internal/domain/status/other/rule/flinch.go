package rule

import (
	"fmt"
	"pob/battle/internal/domain/status"
)

// ひるみ
type Flinch struct {
	kind status.ClearedOnMoveTurnEnd
}

func NewFlinch() *Flinch {
	return &Flinch{kind: status.Flinch}
}

func (f *Flinch) Resolve(ctx status.OtherStatusContext) (bool, bool, string) {
	return true, false, ""
}

func (f *Flinch) Kind() status.OtherCondition { return status.OtherCondition(f.kind) }

func (f *Flinch) Handle(name string) string {
	return fmt.Sprintf("%sはひるんで動けなかった", name)
}
