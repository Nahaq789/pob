package ability

type AbilityId int

type Ability struct {
	originalId   AbilityId
	originalName string
	changedId    *AbilityId
	changedName  *string
	canSwap      bool // 特性を変更できるか
	canCopy      bool // 特性をコピーできるか
}

// 特性交換できない特性たち
var unswappableIds = map[AbilityId]struct{}{
	// おそらく初期リリース時点の第一世代では存在しない
}

// コピーできない特性たち
var uncopyableIds = map[AbilityId]struct{}{
	// おしらく初期リリース時点の第一世代では存在しない
}

func NewAbility(originalId int, originalName string) Ability {
	_, swapOk := unswappableIds[AbilityId(originalId)]
	_, copyOk := uncopyableIds[AbilityId(originalId)]
	return Ability{
		originalId:   AbilityId(originalId),
		originalName: originalName,
		changedId:    nil,
		changedName:  nil,
		canSwap:      !swapOk,
		canCopy:      !copyOk,
	}
}

func (a Ability) GetCurrentId() AbilityId {
	if a.changedId != nil {
		return *a.changedId
	}
	return a.originalId
}

func (a Ability) GetCurrentName() string {
	if a.changedName != nil {
		return *a.changedName
	}

	return a.originalName
}
