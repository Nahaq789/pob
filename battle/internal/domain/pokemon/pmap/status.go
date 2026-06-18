package pmap

import "pob/battle/internal/domain/status"

// その他状態異常(こんらんとか)のハンドラー置き場
var OtherConditionMap = map[status.OtherCondition]OtherStatusHandler{}

type OtherStatusHandler interface {
	// 仮
	Execute(s status.OtherStatus) error
}

// 状態異常のハンドラー置き場
var MainConditionMap = map[status.Condition]MainStatusHandler{}

type MainStatusHandler interface {
	// 仮
	Apply() error
}

// 上書き可能な状態異常の技一覧
var OverridableStatusMoves = map[int]string{
	173: "ねむる",
}
