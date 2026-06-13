package rank

import "fmt"

var basicRank = map[int]float64{
	0:  1.0,
	1:  1.5,
	2:  2.0,
	3:  2.5,
	4:  3.0,
	5:  3.5,
	6:  4.0,
	12: 4.0,
	-1: 0.667,
	-2: 0.5,
	-3: 0.4,
	-4: 0.33,
	-5: 0.286,
	-6: 0.25,
}

const MaxStage = 12
const ResetStage = 0

type BasicRank struct {
	stage int
	value float64
}

func NewBasicRank() BasicRank {
	return BasicRank{stage: 0, value: 1.0}
}

func (r BasicRank) get(n int) (float64, error) {
	f, ok := basicRank[n]
	if !ok {
		return 0.0, fmt.Errorf("unexpected rank: %v", n)
	}
	return f, nil
}

// Up はランクを stage 段階上昇させ、上昇後の BasicRank・メッセージ・エラーを返す。
// stage には通常 1〜6 の変化量を渡す。
// MaxStage を渡すと現在のランクに関わらず強制的に最大まで引き上げる（はらだいこ・いかりのつぼ 等）。
//
// 例:
//
//	r := rank.NewBasicRank()           // stage=0
//	r, msg, _ := r.Up(2)              // → stage=2, msg="が ぐーんとあがった"
//	r, msg, _ = r.Up(rank.MaxStage)   // → stage=6, msg="が さいだいまであがった"
func (r BasicRank) Up(stage int) (BasicRank, string, error) {
	// 「はらだいこ」や「いかりのつぼ」のような
	// 現在の状態に関わらず能力値を最大まであげる場合の処理
	if stage == MaxStage {
		// 取得できないケースはありえないので、エラーチェックなし
		f, _ := r.get(MaxStage)
		return BasicRank{stage: 6, value: f}, message(MaxStage), nil
	}

	var n = r.stage + stage
	if n > 6 {
		n = 6
	}
	f, err := r.get(n)
	if err != nil {
		return BasicRank{}, "", err
	}

	b := BasicRank{stage: n, value: f}

	d := n - r.stage
	return b, message(d), nil
}

// Down はランクを stage 段階下降させ、下降後の BasicRank・メッセージ・エラーを返す。
// stage には通常 1〜6 の変化量（正の値）を渡す。下限は −6 でクランプされる。
//
// 例:
//
//	r := rank.NewBasicRank()      // stage=0
//	r, msg, _ := r.Down(1)       // → stage=-1, msg="が さがった"
//	r, msg, _ = r.Down(3)        // → stage=-4, msg="が がくーんとさがった"
func (r BasicRank) Down(stage int) (BasicRank, string, error) {
	var n = r.stage - stage
	if n < -6 {
		n = -6
	}

	f, err := r.get(n)
	if err != nil {
		return BasicRank{}, "", err
	}

	b := BasicRank{stage: n, value: f}
	d := n - r.stage
	return b, message(d), nil
}

func (r BasicRank) Reset() BasicRank {
	f, _ := r.get(ResetStage)
	return BasicRank{stage: ResetStage, value: f}
}

func message(delta int) string {
	switch {
	case delta == MaxStage:
		return "が さいだいまであがった"
	case delta >= 3:
		return "が ぐぐーんとあがった"
	case delta == 2:
		return "が ぐーんとあがった"
	case delta == 1:
		return "が あがった"
	case delta == -1:
		return "が さがった"
	case delta == -2:
		return "が がくっとさがった"
	case delta <= -3:
		return "が がくーんとさがった"
	default:
		return ""
	}
}
