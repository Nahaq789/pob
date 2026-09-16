package phase

import "pob/battle/internal/domain/damage"

type Phase string

const (
	PhaseMoveResolve Phase = "move_resolve"
	PhasePostMove    Phase = "post_move"
	PhaseEnd         Phase = "end"
)

// ポケモン登場時に発生するイベント（特性）のインターフェース
// 「いかく」等のハンドラーはこちらを実装する
type EntryHandler interface {
	Handle(ctx EntryContext) Result
}

// ポケモン退場時に発生するイベント（特性・技）のインターフェース
// 「すてぜりふ」等のハンドラーはこちらを実装する
type ExitHandler interface {
	Handle(ctx ExitContext) Result
}

// PreMoveHandler は技使用前フェーズのハンドラーのインターフェース。
// ポケモンの行動チェック等を実施する。
type PreMoveHandler interface {
	Handle(ctx PreMoveContext)
}

// DamageModHandler は技解決フェーズの補正値を返すハンドラーのインターフェース。
// 基本ハンドラー（急所・タイプ相性・やけど等）と特殊ハンドラー（技・特性・道具）の両方が実装する。
type DamageModHandler interface {
	Mod(ctx MoveResolveContext) damage.DamageMod
}

// MoveResolveHandler は技の追加効果を適用するインターフェース。
// ダメージ技・変化技・追加効果（状態異常付与・能力ランク変化等）を適用する。
type MoveResolveHandler interface {
	AfterEffect(ctx MoveResolveContext) Result
}

// MoveHandler は技ハンドラーの共通インターフェース。
// 全ての技はこのインターフェースを実装する。
// 補正不要な場合は Mod で空の DamageMod を、追加効果なしの場合は AfterEffect で空の Result を返す。
// 単発技は HitCount で 1 を返す。
type MoveHandler interface {
	DamageModHandler
	MoveResolveHandler
	HitCount() int
}
