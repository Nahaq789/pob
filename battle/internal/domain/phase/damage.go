package phase

import (
	"math/rand/v2"
	"pob/battle/internal/domain/damage"
	"pob/battle/internal/domain/move"
	"pob/battle/internal/domain/pokemon"
)

type DamagePhaseHandler struct {
	registry *Registry
}

func NewDamagePhaseHandler(r *Registry) *DamagePhaseHandler {
	return &DamagePhaseHandler{registry: r}
}

func (d *DamagePhaseHandler) Handle(ctx DamageContext) Result {
	actor := ctx.Battle.PlayerById(ctx.ActorId)
	attacker := actor.Active()

	// TODO
	// ここを混乱判定部分に修正する
	var defender *pokemon.Pokemon
	if ctx.TargetSelf {
		defender = attacker
	} else {
		defender = ctx.Battle.Opponent(actor).Active()
	}

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

	// 基本ハンドラーを順に実行。CritHandler が先頭に登録されている前提で、
	// mod.Crit が確定した時点で ctx.IsCrit を更新し後続ハンドラーが参照できるようにする。
	var mod damage.DamageMod
	for _, h := range d.registry.damageBaseHandlers {
		mod = damage.Merge(mod, h.Mod(ctx))
		if mod.Crit != nil && mod.Crit.Value > 1.0 {
			ctx.IsCrit = true
		}
	}

	// 攻撃側の特性・道具ハンドラー
	if abilityId := int(attacker.Ability().GetCurrentId()); abilityId != 0 {
		if h, ok := d.registry.damageAbilityHandlers[abilityId]; ok {
			mod = damage.Merge(mod, h.Mod(ctx))
		}
	}
	if item := attacker.HeldItem(); item != nil {
		if h, ok := d.registry.damageItemHandlers[int(item.Id())]; ok {
			mod = damage.Merge(mod, h.Mod(ctx))
		}
	}

	// 防御側の特性・道具ハンドラー（自傷時は重複しないようスキップ）
	if attacker != defender {
		if abilityId := int(defender.Ability().GetCurrentId()); abilityId != 0 {
			if h, ok := d.registry.damageAbilityHandlers[abilityId]; ok {
				mod = damage.Merge(mod, h.Mod(ctx))
			}
		}
		if item := defender.HeldItem(); item != nil {
			if h, ok := d.registry.damageItemHandlers[int(item.Id())]; ok {
				mod = damage.Merge(mod, h.Mod(ctx))
			}
		}
	}

	// 技ハンドラー
	if h, ok := d.registry.damageMoveHandlers[ctx.MoveId]; ok {
		mod = damage.Merge(mod, h.Mod(ctx))
	}

	input := damage.NewDamageInput(damage.NewPower(ctx.Power), attack, def, random, mod.ToOptions()...)
	dmg := input.CalcDamage()
	defender.TakeDamage(dmg)

	return Result{
		NextPhase: PhasePostDamage,
	}
}
