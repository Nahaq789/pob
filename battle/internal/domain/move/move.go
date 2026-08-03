package move

import (
	"pob/battle/internal/domain/pp"
	"pob/battle/internal/domain/ptype"
)

type Move struct {
	id          int
	pp          pp.PP
	power       int // 0: 固定・変動ダメージ技
	accuracy    int // 0: 必中
	priority    int
	damageClass DamageClass
	moveType    ptype.Type
}

func NewMove(id int, pp pp.PP, power, accuracy, priority int, damageClass DamageClass, moveType ptype.Type) Move {
	return Move{id: id, pp: pp, power: power, accuracy: accuracy, priority: priority, damageClass: damageClass, moveType: moveType}
}

func (m Move) Id() int                  { return m.id }
func (m Move) PP() pp.PP                { return m.pp }
func (m Move) Power() int               { return m.power }
func (m Move) Accuracy() int            { return m.accuracy }
func (m Move) Priority() int            { return m.priority }
func (m Move) DamageClass() DamageClass { return m.damageClass }
func (m Move) Type() ptype.Type         { return m.moveType }

func (m Move) ConsumePP() Move {
	newPP := m.pp.ConsumeOne()
	return Move{id: m.id, pp: newPP, power: m.power, accuracy: m.accuracy, priority: m.priority, damageClass: m.damageClass, moveType: m.moveType}
}
