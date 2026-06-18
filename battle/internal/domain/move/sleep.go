package move

import (
	"math/rand/v2"
	"pob/battle/internal/domain/vo"
)

// RollSleepCount はねむり状態の継続ターン数を抽選する。
// ねむる技固有のロジックであり、status パッケージは関知しない。
func RollSleepCount() vo.Count {
	return vo.NewCount(2 + rand.IntN(2))
}
