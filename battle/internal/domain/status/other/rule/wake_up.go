package rule

import "pob/battle/internal/domain/status/other"

type ClearedOnWakeUp struct{}

func NewClearedOnWakeUp() *ClearedOnWakeUp {
	return &ClearedOnWakeUp{}
}

func (c *ClearedOnWakeUp) Resolve(ctx other.OtherStatusContext) (bool, bool) {
	// TODO: pokemon.(*Pokemon).MainStatus() が現在コメントアウト中のため未実装。
	// 公開後は以下のロジックで実装する:
	//   player := ctx.Battle.PlayerById(ctx.ActorId)
	//   ms := player.Active().MainStatus()  // *status.MainStatus
	//   if ms == nil || ms.Condition() != status.Sleep {
	//       return true, false
	//   }
	//   return false, false
	_ = ctx
	return false, false
}
