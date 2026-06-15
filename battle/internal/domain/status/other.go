package status

import "pob/battle/internal/domain/vo"

type OtherStatus struct {
	condition OtherCondition
	count     vo.Count
}

var OtherConditionMap = map[OtherCondition]OtherConditionHandler{}
