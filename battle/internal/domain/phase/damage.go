package phase

// DamageContext は「ダメージ計算中」フェーズで handler に渡されるコンテキスト。
// フィールドは今後の設計で詰める。
type DamageContext struct {
	// TODO: actorId, targetId, abilityId, itemId, moveId, move, battle,
	//       power, typeRate, critical, damage 等を追加
	// battle *battle.Battle  // Battle 集約実装後に有効化
}

func NewDamageContext( /* battle *battle.Battle */ ) *DamageContext {
	return &DamageContext{}
}

// DamageHandler は「ダメージ計算中」フェーズで実行される handler の共通 interface。
// やけど時の攻撃半減、天候による威力補正、たいねつ、ふかしのこぶし、連続攻撃などを扱う。
type DamageHandler interface {
	Handle(ctx *DamageContext)
}
