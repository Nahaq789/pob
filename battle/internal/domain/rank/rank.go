package rank

type Rank struct {
	attack, defence, spAttack, spDefence, speed, accuracy, evasion BasicRank
	critical                                                       CriticalRank
}

func NewRank() Rank {
	b := NewBasicRank()
	return Rank{
		attack:    b,
		defence:   b,
		spAttack:  b,
		spDefence: b,
		speed:     b,
		accuracy:  b,
		evasion:   b,
		critical:  NewCriticalRank(),
	}
}

func (r Rank) Attack() BasicRank      { return r.attack }
func (r Rank) Defence() BasicRank     { return r.defence }
func (r Rank) SpAttack() BasicRank    { return r.spAttack }
func (r Rank) SpDefence() BasicRank   { return r.spDefence }
func (r Rank) Speed() BasicRank       { return r.speed }
func (r Rank) Critical() CriticalRank { return r.critical }
func (r Rank) Accuracy() BasicRank    { return r.accuracy }
func (r Rank) Evasion() BasicRank     { return r.evasion }

// 自分の命中ランクと相手の回避ランクの差分から命中補正値を返す。
func (r Rank) AccuracyRank(oe BasicRank) (AccuracyRank, error) {
	a, err := NewAccuracyRank(r.accuracy.stage, oe.stage)
	if err != nil {
		return AccuracyRank{}, err
	}

	return a, nil
}
