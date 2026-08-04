package phase

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

type PreDamageHandler interface {
	Handle(ctx string)
}
