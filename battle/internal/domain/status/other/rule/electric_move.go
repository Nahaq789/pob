package rule

import "pob/battle/internal/domain/status/other"

type ClearedOnElectricMove struct{}

func NewClearedOnElectricMove() *ClearedOnElectricMove {
	return &ClearedOnElectricMove{}
}

func (c *ClearedOnElectricMove) Resolve(ctx other.OtherStatusContext) (bool, bool) {
	// TODO: でんきタイプ判定には ptype.Electric との比較が必要だが、
	// 依存制約(status/other は battle と status のみに依存可)により ptype を import できない。
	// 以下いずれかで解消後に実装する:
	//   (a) OtherStatusContext に MoveType ptype.Type フィールドを追加する
	//   (b) battle パッケージに「技IDからタイプを返すヘルパー」を追加する
	//   (c) 依存制約を ptype まで緩和する
	// 実装イメージ(a案):
	//   if ctx.MoveType == ptype.Electric {
	//       return true, false
	//   }
	//   return false, false
	_ = ctx
	return false, false
}
