package nature

import (
	"pob/battle/internal/domain/item"
	"pob/pkg/stats"
)

type Nature struct {
	name  string
	stats stats.NatureModifier
}

func NewNature(n string) (Nature, error) {
	mod, err := stats.GetNatureModifiers(n)
	if err != nil {
		return Nature{}, err
	}
	return Nature{name: n, stats: mod}, nil
}

func (n Nature) DislikedConfuseFlavor() item.Flavor {
	switch {
	case n.stats.A < 1:
		return item.FlavorSpicy
	case n.stats.B < 1:
		return item.FlavorSour
	case n.stats.C < 1:
		return item.FlavorDry
	case n.stats.D < 1:
		return item.FlavorBitter
	case n.stats.S < 1:
		return item.FlavorSweet
	default:
		return "" // 無補正性格、嫌いな味なし
	}
}
