package phase

type Phase string

const (
	PhaseDamage     Phase = "damage"
	PhasePostDamage Phase = "post_damage"
	PhaseEnd        Phase = "end"
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

// ダメージ計算前のハンドラー
// ポケモンの行動チェックや素早さ順の解決等実施する
type PreDamageHandler interface {
	Handle(ctx string)
}
