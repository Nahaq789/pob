package move

import "pob/battle/internal/domain/pp"

type Move struct {
	id int
	pp pp.PP
}

func NewMove(id int, pp pp.PP) Move {
	return Move{id: id, pp: pp}
}

func (m Move) Id() int   { return m.id }
func (m Move) PP() pp.PP { return m.pp }

func (m Move) ConsumePP() Move {
	newPP := m.pp.ConsumeOne()
	return Move{id: m.id, pp: newPP}
}
