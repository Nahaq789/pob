package phase

// MoveSelectContext は「技選択前」フェーズで handler に渡されるコンテキスト。
// フィールドは今後の設計で詰める。
type MoveSelectContext struct {
	// TODO: actorId, abilityId, itemId, battle, availableMoves, canAct 等を追加
	// battle *battle.Battle  // Battle 集約実装後に有効化
}

func NewMoveSelectContext( /* battle *battle.Battle */ ) *MoveSelectContext {
	return &MoveSelectContext{}
}

// MoveSelectHandler は「技選択前」フェーズで実行される handler の共通 interface。
// こだわり系道具、ねむり/こおり等による行動不能判定などを扱う。
type MoveSelectHandler interface {
	Handle(ctx *MoveSelectContext)
}
