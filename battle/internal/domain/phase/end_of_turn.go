package phase

// EffectKind はターン終了フェーズで発動する効果の種別を一意に識別する。
// handler 同士の競合調停（例: ポイズンヒール vs 毒ダメージ）に使う。
type EffectKind string

const (
	EffectPoisonTick    EffectKind = "poison_tick"
	EffectBurnTick      EffectKind = "burn_tick"
	EffectSandstormTick EffectKind = "sandstorm_tick"
	EffectHailTick      EffectKind = "hail_tick"
	// 必要に応じて追加
)

// EndOfTurnContext はターン終了フェーズで全 handler に渡されるコンテキスト。
// ctx は1つだけ生成され、両プレイヤー分の handled 状態を保持する。
type EndOfTurnContext struct {
	// battle *battle.Battle  // Battle 集約への参照。Battle 実装後に有効化
	handled map[string]map[EffectKind]bool // actorId -> effects
}

func NewEndOfTurnContext( /* battle *battle.Battle */ ) *EndOfTurnContext {
	return &EndOfTurnContext{
		handled: make(map[string]map[EffectKind]bool),
	}
}

func (c *EndOfTurnContext) MarkHandled(actorId string, e EffectKind) {
	if c.handled[actorId] == nil {
		c.handled[actorId] = make(map[EffectKind]bool)
	}
	c.handled[actorId][e] = true
}

func (c *EndOfTurnContext) IsHandled(actorId string, e EffectKind) bool {
	if c.handled[actorId] == nil {
		return false
	}
	return c.handled[actorId][e]
}

// EndOfTurnHandler はターン終了フェーズで実行される handler の共通 interface。
// status/rule・ability・ground/rule の全 handler がこの interface を満たし、
// registry にフラットに登録される。
//
// 各 handler は1回だけ呼ばれる（actor数に依らない）。
// 「片方処理か、両方処理か、actor 非依存か」は handler 自身が判断する。
type EndOfTurnHandler interface {
	Handle(ctx *EndOfTurnContext)
}
