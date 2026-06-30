package phase

// EntryContext は「場に出たとき」フェーズで handler に渡されるコンテキスト。
// フィールドは今後の設計で詰める。
type EntryContext struct {
	// TODO: actorId, abilityId, itemId, battle 等を必要に応じて追加
	// battle *battle.Battle  // Battle 集約実装後に有効化
}

func NewEntryContext( /* battle *battle.Battle */ ) *EntryContext {
	return &EntryContext{}
}

// EntryHandler は「場に出たとき」フェーズで実行される handler の共通 interface。
// 特性発動(いかく・てんねん等)、天候変化系特性、入場時アイテム効果などを扱う。
type EntryHandler interface {
	Handle(ctx *EntryContext)
}
