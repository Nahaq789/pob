package rank

import "fmt"

var accuracyRankTable = map[int][2]int{
	0:  {3, 3},
	1:  {4, 3},
	2:  {5, 3},
	3:  {6, 3},
	4:  {7, 3},
	5:  {8, 3},
	6:  {9, 3},
	-1: {3, 4},
	-2: {3, 5},
	-3: {3, 6},
	-4: {3, 7},
	-5: {3, 8},
	-6: {3, 9},
}

type AccuracyRank struct {
	stage int
	value [2]int
}

func NewAccuracyRank(h, e int) (AccuracyRank, error) {
	s := min(max(h-e, -6), 6)
	f, ok := accuracyRankTable[s]
	if !ok {
		return AccuracyRank{}, fmt.Errorf("accuracy rank: invalid stage %d", s)
	}

	return AccuracyRank{stage: s, value: f}, nil
}

func (r AccuracyRank) Value() [2]int {
	return r.value
}
