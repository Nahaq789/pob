package phase

// PostMoveContext は「ダメージ計算後」フェーズで handler に渡されるコンテキスト。
// フィールドは今後の設計で詰める。
type PostMoveContext struct {
	// TODO: actorId, targetId, abilityId, itemId, moveId, move, battle, dealtDamage 等を追加
	// battle *battle.Battle  // Battle 集約実装後に有効化
}

func NewPostMoveContext( /* battle *battle.Battle */ ) *PostMoveContext {
	return &PostMoveContext{}
}

// PostMoveHandler は「ダメージ計算後」フェーズで実行される handler の共通 interface。
// じきゅうりょく、追加効果(状態異常付与等)、変化技の効果適用などを扱う。
type PostMoveHandler interface {
	Handle(ctx *PostMoveContext)
}
