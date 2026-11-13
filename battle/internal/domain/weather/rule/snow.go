package rule

import (
	"pob/battle/internal/domain/damage"
	"pob/battle/internal/domain/phase"
)

var _ phase.StatModHandler = (*SnowHandler)(nil)

type SnowHandler struct{}

func NewSnowHandler() *SnowHandler { return &SnowHandler{} }

func (s *SnowHandler) StatMod(ctx phase.MoveResolveContext) damage.StatMod {
	return damage.StatMod{}
}
