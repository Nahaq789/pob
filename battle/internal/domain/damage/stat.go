package damage

type StatMod struct {
	AttackerStat *ModValue
	DefenderStat *ModValue
}

func MergeStat(base, incoming StatMod) StatMod {
	return StatMod{
		AttackerStat: mergeValue(base.AttackerStat, incoming.AttackerStat),
		DefenderStat: mergeValue(base.DefenderStat, incoming.DefenderStat),
	}
}

func (s StatMod) Apply(attack, defense int) (int, int) {
	if s.AttackerStat != nil {
		attack = roundHalfDown(float64(attack) * s.AttackerStat.Value)
	}
	if s.DefenderStat != nil {
		defense = roundHalfDown(float64(defense) * s.DefenderStat.Value)
	}
	return attack, defense
}
