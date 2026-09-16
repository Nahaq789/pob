package phase

import (
	"math/rand/v2"
	"pob/battle/internal/domain/damage"
	"pob/battle/internal/domain/move"
)

type MoveResolvePhaseHandler struct {
	registry *Registry
}

func NewMoveResolvePhaseHandler(r *Registry) *MoveResolvePhaseHandler {
	return &MoveResolvePhaseHandler{registry: r}
}

func (d *MoveResolvePhaseHandler) Handle(ctx MoveResolveContext) Result {
	actor := ctx.Battle.PlayerById(ctx.ActorId)
	attacker := actor.Active()
	defender := ctx.Battle.Opponent(actor).Active()

	var attack, def int
	switch ctx.Category {
	case move.DamageClassPhysical:
		attack = attacker.AttackStat()
		def = defender.DefenseStat()
	default:
		attack = attacker.SpAttackStat()
		def = defender.SpDefenseStat()
	}

	random := 85 + rand.IntN(16)

	// 混乱時の処理
	if ctx.TargetSelf {
		input := damage.NewDamageInput(damage.NewPower(ctx.Power), attack, def, random)
		dmg := input.CalcDamage()
		attacker.TakeDamage(dmg)

		return Result{
			Messages:  []string{"わけもわからず自分を攻撃した"},
			NextPhase: PhasePostMove,
		}
	}

	// 基本ハンドラーを順に実行。CritHandler が先頭に登録されている前提で、
	// mod.Crit が確定した時点で ctx.IsCrit を更新し後続ハンドラーが参照できるようにする。
	var mod damage.DamageMod
	for _, h := range d.registry.resolveBaseHandlers {
		mod = damage.Merge(mod, h.Mod(ctx))
		if mod.Crit != nil && mod.Crit.Value > 1.0 {
			ctx.IsCrit = true
		}
	}

	// 攻撃側の特性・道具（スナイパー・いのちのたま等）
	if abilityId := int(attacker.Ability().GetCurrentId()); abilityId != 0 {
		if h, ok := d.registry.resolveAbilityHandlers[abilityId]; ok {
			mod = damage.Merge(mod, h.Mod(ctx))
		}
	}
	if item := attacker.HeldItem(); item != nil {
		if h, ok := d.registry.resolveItemHandlers[int(item.Id())]; ok {
			mod = damage.Merge(mod, h.Mod(ctx))
		}
	}

	// 防御側の特性・道具（マルチスケイル・フィルター・もふもふ等）
	if abilityId := int(defender.Ability().GetCurrentId()); abilityId != 0 {
		if h, ok := d.registry.resolveAbilityHandlers[abilityId]; ok {
			mod = damage.Merge(mod, h.Mod(ctx))
		}
	}
	if item := defender.HeldItem(); item != nil {
		if h, ok := d.registry.resolveItemHandlers[int(item.Id())]; ok {
			mod = damage.Merge(mod, h.Mod(ctx))
		}
	}
	if h, ok := d.registry.resolveMoveHandlers[ctx.MoveId]; ok {
		mod = damage.Merge(mod, h.Mod(ctx))
	}

	input := damage.NewDamageInput(damage.NewPower(ctx.Power), attack, def, random, mod.ToOptions()...)
	dmg := input.CalcDamage()
	defender.TakeDamage(dmg)

	return Result{
		NextPhase: PhasePostMove,
	}
}
