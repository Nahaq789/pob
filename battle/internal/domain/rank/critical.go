package rank

import "fmt"

var criticalRankMap = map[int]float64{
	0: 0.417,
	1: 0.125,
	2: 0.5,
	3: 1.0,
}

type CriticalRank struct {
	stage int
	value float64
}

func NewCriticalRank() CriticalRank {
	f, _ := criticalRankMap[ResetStage]
	return CriticalRank{stage: ResetStage, value: f}
}

func (r CriticalRank) get(n int) (float64, error) {
	f, ok := criticalRankMap[n]
	if !ok {
		return 0.0, fmt.Errorf("")
	}

	return f, nil
}

func (r CriticalRank) Up(stage int) (CriticalRank, error) {
	var n = r.stage + stage
	if n > CriticalMaxStage {
		n = 3
	}

	f, err := r.get(n)
	if err != nil {
		return CriticalRank{}, err
	}

	return CriticalRank{stage: n, value: f}, nil
}

func (r CriticalRank) Down(stage int) (CriticalRank, error) {
	var n = r.stage - stage
	if n < ResetStage {
		n = ResetStage
	}

	f, err := r.get(n)
	if err != nil {
		return CriticalRank{}, err
	}

	return CriticalRank{stage: n, value: f}, nil
}

func (r CriticalRank) Reset() CriticalRank {
	return NewCriticalRank()
}
