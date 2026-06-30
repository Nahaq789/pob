package phase

// PreDamageContext は「ダメージ計算前」フェーズで handler に渡されるコンテキスト。
// フィールドは今後の設計で詰める。
type PreDamageContext struct {
	// TODO: actorId, targetId, abilityId, itemId, moveId, move, battle 等を追加
	// battle *battle.Battle  // Battle 集約実装後に有効化
}

func NewPreDamageContext( /* battle *battle.Battle */ ) *PreDamageContext {
	return &PreDamageContext{}
}

// PreDamageHandler は「ダメージ計算前」フェーズで実行される handler の共通 interface。
// タイプ変更系特性、ジュエル系道具、溜め技処理などを扱う。
type PreDamageHandler interface {
	Handle(ctx *PreDamageContext)
}
