package status

import "pob/battle/internal/domain/ptype"

// OtherStatusContext は OtherStatus.Resolve に渡す最小限の情報。
// EntryContext / PreDamageContext と同じ設計思想。
type OtherStatusContext struct {
	// ActorId は行動したポケモンの識別子
	ActorId string
	// ActorName はメッセージ生成に使用するポケモン名
	ActorName string
	// MoveId は使用した技のID
	MoveId string
	// ActionSucceeded は技が成功したかどうか
	ActionSucceeded bool
	// MoveType は使用した技のタイプ。ClearedOnElectricMove の判定に使用する
	MoveType ptype.Type
	// MainCondition は現在のメイン状態異常。ClearedOnWakeUp の判定に使用する（None なら状態異常なし）
	MainCondition Condition
}

// OtherStatus は「わざを使ったポケモンに発生し、行動選択を制約する状態変化」の共通インターフェース。
type OtherStatus interface {
	// Resolve はそのポケモンの1ターン分の行動解決後に呼ばれる。
	// 戻り値:
	//   cleared      : true の場合、呼び出し側はこの OtherStatus をポケモンから除去し、
	//                   技選択(MoveSelect)を通常通り可能な状態に戻す
	//   addConfusion : true の場合、呼び出し側は cleared 除去の直後に status.Confusion を付与する
	//                  (こんらん付与自体はOrchestrator側の責務。Resolve内でSetStatusは呼ばない)
	//   message      : 状態異常解除時など、表示すべきメッセージ。ない場合は空文字
	Resolve(ctx OtherStatusContext) (cleared bool, addConfusion bool, message string)
	Kind() OtherCondition
}
