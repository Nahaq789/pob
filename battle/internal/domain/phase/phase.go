package phase

// ポケモン登場時に発生するイベント（特性）のインターフェース
// 「いかく」等のハンドラーはこちらを実装する
type EntryHandler interface {
	Handle(ctx EntryContext) Result
}
