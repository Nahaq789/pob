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

// MoveResolveHandler は技解決フェーズの補正値を返すハンドラーのインターフェース。
// 基本ハンドラー（急所・タイプ相性・やけど等）と特殊ハンドラー（技・特性・道具）の両方が実装する。
type MoveResolveHandler interface {
	Mod(ctx MoveResolveContext) damage.DamageMod
}
