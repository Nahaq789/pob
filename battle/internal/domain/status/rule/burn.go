package rule

import (
	"pob/battle/internal/domain/damage"
	"pob/battle/internal/domain/move"
	"pob/battle/internal/domain/phase"
)

var _ phase.DamageModHandler = (*BurnHandler)(nil)

type BurnHandler struct{}

func NewBurnHandler() *BurnHandler { return &BurnHandler{} }

func (b *BurnHandler) Mod(ctx phase.MoveResolveContext) damage.DamageMod {
	if ctx.Category != move.DamageClassPhysical {
		return damage.DamageMod{}
	}

	attacker := ctx.Battle.PlayerById(ctx.ActorId).Active()
	if s := attacker.Status().Main(); s == nil || !s.IsBurn() {
		return damage.DamageMod{}
	}

	return damage.DamageMod{Burn: damage.Multiply(0.5)}
}
