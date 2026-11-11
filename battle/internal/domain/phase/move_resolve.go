package phase

import (
	"math/rand/v2"
	"pob/battle/internal/domain/damage"
	"pob/battle/internal/domain/move"
	"pob/battle/internal/domain/pokemon"
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

	var messages []string

	moveHandler, hasMoveHandler := d.registry.resolveMoveHandlers[ctx.MoveId]

	if ctx.Category != move.DamageClassStatus {
		hitCount := 1
		if hasMoveHandler {
			hitCount = moveHandler.HitCount()
		}

		for range hitCount {
			// 固定ダメージ技は急所判定を行わないため、分岐前にリセットして前の発の結果を持ち越さない。
			ctx.IsCrit = false

			var dmg int
			if f, ok := moveHandler.(FixedDamager); ok {
				dmg = f.FixedDamage(ctx)
			} else {
				dmg = d.calcDamage(&ctx, attacker, defender, attack, def, moveHandler)
			}
			defender.TakeDamage(dmg)

			// 被ダメージ時発動（じきゅうりょく・きのみ等）
			if abilityId := int(defender.Ability().GetCurrentId()); abilityId != 0 {
				if h, ok := d.registry.onHitAbilityHandlers[abilityId]; ok {
					messages = append(messages, h.OnHit(ctx).Messages...)
				}
			}
			if item := defender.HeldItem(); item != nil {
				if h, ok := d.registry.onHitItemHandlers[int(item.Id())]; ok {
					messages = append(messages, h.OnHit(ctx).Messages...)
				}
			}
		}
	}

	// 追加効果計算
	if hasMoveHandler {
		messages = append(messages, moveHandler.AfterEffect(ctx).Messages...)
	}

	return Result{
		Messages:  messages,
		NextPhase: PhasePostMove,
	}
}

// calcDamage は各ハンドラーの補正値を Merge し、ダメージ計算式で 1 発分のダメージを返す。
// 能力値の補正（StatMod）は攻撃・防御の実数値に適用してから計算式に渡す。
// ctx.IsCrit は後続の被弾時ハンドラーも参照するため、ポインタで受け取り更新する。
func (d *MoveResolvePhaseHandler) calcDamage(ctx *MoveResolveContext, attacker, defender *pokemon.Pokemon, attack, def int, moveHandler MoveHandler) int {
	random := 85 + rand.IntN(16)

	var mod damage.DamageMod
	var statMod damage.StatMod // 能力値の補正（すなあらし・こだわりハチマキ等）

	// 基本ハンドラーを順に実行。CritHandler が先頭に登録されている前提で、
	// mod.Crit が確定した時点で ctx.IsCrit を更新し後続ハンドラーが参照できるようにする。
	for _, h := range d.registry.resolveBaseHandlers {
		mod = damage.Merge(mod, h.Mod(*ctx))
		if mod.Crit != nil && mod.Crit.Value > 1.0 {
			ctx.IsCrit = true
		}
	}
	for _, h := range d.registry.statModBaseHandlers {
		statMod = damage.MergeStat(statMod, h.StatMod(*ctx))
	}

	// 攻撃側の特性・道具（スナイパー・いのちのたま・ちからもち・こだわりハチマキ等）
	if abilityId := int(attacker.Ability().GetCurrentId()); abilityId != 0 {
		if h, ok := d.registry.resolveAbilityHandlers[abilityId]; ok {
			mod = damage.Merge(mod, h.Mod(*ctx))
		}
		if h, ok := d.registry.statModAbilityHandlers[abilityId]; ok {
			statMod = damage.MergeStat(statMod, h.StatMod(*ctx))
		}
	}
	if item := attacker.HeldItem(); item != nil {
		if h, ok := d.registry.resolveItemHandlers[int(item.Id())]; ok {
			mod = damage.Merge(mod, h.Mod(*ctx))
		}
		if h, ok := d.registry.statModItemHandlers[int(item.Id())]; ok {
			statMod = damage.MergeStat(statMod, h.StatMod(*ctx))
		}
	}

	// 防御側の特性・道具（マルチスケイル・フィルター・もふもふ・ふしぎなうろこ・とつげきチョッキ等）
	if abilityId := int(defender.Ability().GetCurrentId()); abilityId != 0 {
		if h, ok := d.registry.resolveAbilityHandlers[abilityId]; ok {
			mod = damage.Merge(mod, h.Mod(*ctx))
		}
		if h, ok := d.registry.statModAbilityHandlers[abilityId]; ok {
			statMod = damage.MergeStat(statMod, h.StatMod(*ctx))
		}
	}
	if item := defender.HeldItem(); item != nil {
		if h, ok := d.registry.resolveItemHandlers[int(item.Id())]; ok {
			mod = damage.Merge(mod, h.Mod(*ctx))
		}
		if h, ok := d.registry.statModItemHandlers[int(item.Id())]; ok {
			statMod = damage.MergeStat(statMod, h.StatMod(*ctx))
		}
	}
	if moveHandler != nil {
		mod = damage.Merge(mod, moveHandler.Mod(*ctx))
	}

	attack, def = statMod.Apply(attack, def)
	input := damage.NewDamageInput(damage.NewPower(ctx.Power), attack, def, random, mod.ToOptions()...)
	return input.CalcDamage()
}
