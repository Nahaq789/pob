package phase_test

import (
	"strings"
	"testing"

	"pob/battle/internal/domain/ability"
	"pob/battle/internal/domain/battle"
	"pob/battle/internal/domain/hp"
	"pob/battle/internal/domain/move"
	"pob/battle/internal/domain/nature"
	"pob/battle/internal/domain/phase"
	"pob/battle/internal/domain/player"
	"pob/battle/internal/domain/pokemon"
	"pob/battle/internal/domain/pp"
	"pob/battle/internal/domain/ptype"
	"pob/battle/internal/domain/rank"
	"pob/battle/internal/domain/status"
	"pob/battle/internal/domain/status/other/rule"
	"pob/battle/internal/domain/vo"
)

// ── helpers ───────────────────────────────────────────────────────────────────

func newPreMovePokemon(name string, st status.Status, moveIds ...int) *pokemon.Pokemon {
	a := ability.NewAbility(0, "")
	var moves [4]*move.Move
	for i, id := range moveIds {
		if i >= 4 {
			break
		}
		m := move.NewMove(id, pp.NewPP(vo.NewCount(10), 10), 0, 0, 0, move.DamageClassStatus, ptype.Normal)
		moves[i] = &m
	}
	return pokemon.NewPokemon(
		1, name, 0, "",
		[2]ptype.Type{ptype.Normal, ptype.Normal},
		pokemon.BaseStats{},
		pokemon.RealStats{Speed: 100},
		nature.Nature{},
		&a,
		moves,
		hp.NewHP(100),
		rank.NewRank(),
		st,
		nil, nil,
		false,
	)
}

// newPreMoveBattle は actor を p1 のアクティブポケモンとしてバトルを構築する。
// 返値の actorId は PreMoveContext に渡す p1 のプレイヤー ID。
func newPreMoveBattle(actor *pokemon.Pokemon) (*battle.Battle, string) {
	dummy1 := newPreMovePokemon("ダミー1", status.NewStatus(), 1)
	dummy2 := newPreMovePokemon("ダミー2", status.NewStatus(), 1)
	opp1 := newPreMovePokemon("あいて1", status.NewStatus(), 1)
	opp2 := newPreMovePokemon("あいて2", status.NewStatus(), 1)
	opp3 := newPreMovePokemon("あいて3", status.NewStatus(), 1)

	p1 := player.NewPlayer("actor", "Player1", [6]*pokemon.Pokemon{actor, dummy1, dummy2}, nil)
	p2 := player.NewPlayer("opp", "Player2", [6]*pokemon.Pokemon{opp1, opp2, opp3}, nil)

	if err := p1.Select([3]int{0, 1, 2}); err != nil {
		panic(err)
	}
	if err := p2.Select([3]int{0, 1, 2}); err != nil {
		panic(err)
	}

	return battle.NewBattle(p1, p2), "actor"
}

func containsMsg(msgs []string, sub string) bool {
	for _, m := range msgs {
		if strings.Contains(m, sub) {
			return true
		}
	}
	return false
}

// ── 1. 状態異常なし ────────────────────────────────────────────────────────────

func TestPreMoveHandle_NoStatus(t *testing.T) {
	t.Run("状態異常なし: PhaseMoveResolve に進む", func(t *testing.T) {
		actor := newPreMovePokemon("ピカチュウ", status.NewStatus(), 1)
		b, actorId := newPreMoveBattle(actor)

		result := phase.NewPreMovePhaseHandler().Handle(phase.NewPreMoveContext(actorId, 1, b))

		if result.NextPhase != phase.PhaseMoveResolve {
			t.Errorf("expected PhaseMoveResolve, got %v", result.NextPhase)
		}
	})

	t.Run("状態異常なし: メッセージなし", func(t *testing.T) {
		actor := newPreMovePokemon("ピカチュウ", status.NewStatus(), 1)
		b, actorId := newPreMoveBattle(actor)

		result := phase.NewPreMovePhaseHandler().Handle(phase.NewPreMoveContext(actorId, 1, b))

		if len(result.Messages) != 0 {
			t.Errorf("expected no messages, got %v", result.Messages)
		}
	})
}

// ── 2. ねむり ─────────────────────────────────────────────────────────────────

func TestPreMoveHandle_Sleep(t *testing.T) {
	t.Run("眠り継続中: PhaseEnd", func(t *testing.T) {
		ms, err := status.NewSleep(vo.NewCount(2))
		if err != nil {
			t.Fatal(err)
		}
		st := status.NewStatusWith(&ms, nil)
		actor := newPreMovePokemon("カビゴン", st, 1)
		b, actorId := newPreMoveBattle(actor)

		result := phase.NewPreMovePhaseHandler().Handle(phase.NewPreMoveContext(actorId, 1, b))

		if result.NextPhase != phase.PhaseEnd {
			t.Errorf("expected PhaseEnd, got %v", result.NextPhase)
		}
	})

	t.Run("眠り継続中: 睡眠メッセージあり", func(t *testing.T) {
		ms, err := status.NewSleep(vo.NewCount(2))
		if err != nil {
			t.Fatal(err)
		}
		st := status.NewStatusWith(&ms, nil)
		actor := newPreMovePokemon("カビゴン", st, 1)
		b, actorId := newPreMoveBattle(actor)

		result := phase.NewPreMovePhaseHandler().Handle(phase.NewPreMoveContext(actorId, 1, b))

		if !containsMsg(result.Messages, "ぐうぐう眠っている") {
			t.Errorf("expected sleep message, got %v", result.Messages)
		}
	})
}

// ── 3. こおり ─────────────────────────────────────────────────────────────────

// NewMainStatus(Freeze) は count=0 で生成されるため IsFrozen が即時 false を返し、解凍扱いとなる。
func TestPreMoveHandle_Freeze(t *testing.T) {
	t.Run("解凍: PhaseMoveResolve に進む", func(t *testing.T) {
		ms, err := status.NewMainStatus(status.Freeze)
		if err != nil {
			t.Fatal(err)
		}
		st := status.NewStatusWith(&ms, nil)
		actor := newPreMovePokemon("ラプラス", st, 1)
		b, actorId := newPreMoveBattle(actor)

		result := phase.NewPreMovePhaseHandler().Handle(phase.NewPreMoveContext(actorId, 1, b))

		if result.NextPhase != phase.PhaseMoveResolve {
			t.Errorf("expected PhaseMoveResolve, got %v", result.NextPhase)
		}
	})

	t.Run("解凍: 解凍メッセージあり", func(t *testing.T) {
		ms, err := status.NewMainStatus(status.Freeze)
		if err != nil {
			t.Fatal(err)
		}
		st := status.NewStatusWith(&ms, nil)
		actor := newPreMovePokemon("ラプラス", st, 1)
		b, actorId := newPreMoveBattle(actor)

		result := phase.NewPreMovePhaseHandler().Handle(phase.NewPreMoveContext(actorId, 1, b))

		if !containsMsg(result.Messages, "こおりが溶けた") {
			t.Errorf("expected thaw message, got %v", result.Messages)
		}
	})
}

// ── 4. ひるみ ─────────────────────────────────────────────────────────────────

func TestPreMoveHandle_Flinch(t *testing.T) {
	t.Run("ひるみあり: PhaseEnd", func(t *testing.T) {
		fl := rule.NewFlinch()
		st := status.NewStatusWith(nil, []status.OtherStatus{fl})
		actor := newPreMovePokemon("ピカチュウ", st, 1)
		b, actorId := newPreMoveBattle(actor)

		result := phase.NewPreMovePhaseHandler().Handle(phase.NewPreMoveContext(actorId, 1, b))

		if result.NextPhase != phase.PhaseEnd {
			t.Errorf("expected PhaseEnd, got %v", result.NextPhase)
		}
	})

	t.Run("ひるみあり: ひるみメッセージあり", func(t *testing.T) {
		fl := rule.NewFlinch()
		st := status.NewStatusWith(nil, []status.OtherStatus{fl})
		actor := newPreMovePokemon("ピカチュウ", st, 1)
		b, actorId := newPreMoveBattle(actor)

		result := phase.NewPreMovePhaseHandler().Handle(phase.NewPreMoveContext(actorId, 1, b))

		if !containsMsg(result.Messages, "ひるんで") {
			t.Errorf("expected flinch message, got %v", result.Messages)
		}
	})
}

// ── 5. アンコール ─────────────────────────────────────────────────────────────

func TestPreMoveHandle_Encore(t *testing.T) {
	t.Run("アンコール技以外を選択: PhaseEnd", func(t *testing.T) {
		enc := rule.NewEncore(1)
		st := status.NewStatusWith(nil, []status.OtherStatus{enc})
		actor := newPreMovePokemon("ピクシー", st, 1, 2)
		b, actorId := newPreMoveBattle(actor)

		result := phase.NewPreMovePhaseHandler().Handle(phase.NewPreMoveContext(actorId, 2, b))

		if result.NextPhase != phase.PhaseEnd {
			t.Errorf("expected PhaseEnd, got %v", result.NextPhase)
		}
	})

	t.Run("アンコール技以外を選択: アンコールメッセージあり", func(t *testing.T) {
		enc := rule.NewEncore(1)
		st := status.NewStatusWith(nil, []status.OtherStatus{enc})
		actor := newPreMovePokemon("ピクシー", st, 1, 2)
		b, actorId := newPreMoveBattle(actor)

		result := phase.NewPreMovePhaseHandler().Handle(phase.NewPreMoveContext(actorId, 2, b))

		if !containsMsg(result.Messages, "アンコール") {
			t.Errorf("expected encore message, got %v", result.Messages)
		}
	})

	t.Run("アンコール技を選択: PhaseMoveResolve に進む", func(t *testing.T) {
		enc := rule.NewEncore(1)
		st := status.NewStatusWith(nil, []status.OtherStatus{enc})
		actor := newPreMovePokemon("ピクシー", st, 1, 2)
		b, actorId := newPreMoveBattle(actor)

		result := phase.NewPreMovePhaseHandler().Handle(phase.NewPreMoveContext(actorId, 1, b))

		if result.NextPhase != phase.PhaseMoveResolve {
			t.Errorf("expected PhaseMoveResolve, got %v", result.NextPhase)
		}
	})

	t.Run("アンコール技を選択: メッセージなし", func(t *testing.T) {
		enc := rule.NewEncore(1)
		st := status.NewStatusWith(nil, []status.OtherStatus{enc})
		actor := newPreMovePokemon("ピクシー", st, 1, 2)
		b, actorId := newPreMoveBattle(actor)

		result := phase.NewPreMovePhaseHandler().Handle(phase.NewPreMoveContext(actorId, 1, b))

		if len(result.Messages) != 0 {
			t.Errorf("expected no messages, got %v", result.Messages)
		}
	})
}

// ── 6. かなしばり ─────────────────────────────────────────────────────────────

func TestPreMoveHandle_MoveDisabled(t *testing.T) {
	t.Run("封じられた技を選択: PhaseEnd", func(t *testing.T) {
		md := rule.NewMoveDisabled(4, 1)
		st := status.NewStatusWith(nil, []status.OtherStatus{md})
		actor := newPreMovePokemon("ゲンガー", st, 1, 2)
		b, actorId := newPreMoveBattle(actor)

		result := phase.NewPreMovePhaseHandler().Handle(phase.NewPreMoveContext(actorId, 1, b))

		if result.NextPhase != phase.PhaseEnd {
			t.Errorf("expected PhaseEnd, got %v", result.NextPhase)
		}
	})

	t.Run("封じられた技を選択: かなしばりメッセージあり", func(t *testing.T) {
		md := rule.NewMoveDisabled(4, 1)
		st := status.NewStatusWith(nil, []status.OtherStatus{md})
		actor := newPreMovePokemon("ゲンガー", st, 1, 2)
		b, actorId := newPreMoveBattle(actor)

		result := phase.NewPreMovePhaseHandler().Handle(phase.NewPreMoveContext(actorId, 1, b))

		if !containsMsg(result.Messages, "かなしばり") {
			t.Errorf("expected move-disabled message, got %v", result.Messages)
		}
	})

	t.Run("封じられていない技を選択: PhaseMoveResolve に進む", func(t *testing.T) {
		md := rule.NewMoveDisabled(4, 1)
		st := status.NewStatusWith(nil, []status.OtherStatus{md})
		actor := newPreMovePokemon("ゲンガー", st, 1, 2)
		b, actorId := newPreMoveBattle(actor)

		result := phase.NewPreMovePhaseHandler().Handle(phase.NewPreMoveContext(actorId, 2, b))

		if result.NextPhase != phase.PhaseMoveResolve {
			t.Errorf("expected PhaseMoveResolve, got %v", result.NextPhase)
		}
	})
}

// ── 7. こんらん ───────────────────────────────────────────────────────────────

func TestPreMoveHandle_Confusion(t *testing.T) {
	t.Run("残り1ターン: 今ターンで解除 → PhaseMoveResolve", func(t *testing.T) {
		cf := rule.NewConfusion(1)
		st := status.NewStatusWith(nil, []status.OtherStatus{cf})
		actor := newPreMovePokemon("フシギダネ", st, 1)
		b, actorId := newPreMoveBattle(actor)

		result := phase.NewPreMovePhaseHandler().Handle(phase.NewPreMoveContext(actorId, 1, b))

		if result.NextPhase != phase.PhaseMoveResolve {
			t.Errorf("expected PhaseMoveResolve, got %v", result.NextPhase)
		}
	})

	t.Run("残り1ターン: 解除メッセージあり", func(t *testing.T) {
		cf := rule.NewConfusion(1)
		st := status.NewStatusWith(nil, []status.OtherStatus{cf})
		actor := newPreMovePokemon("フシギダネ", st, 1)
		b, actorId := newPreMoveBattle(actor)

		result := phase.NewPreMovePhaseHandler().Handle(phase.NewPreMoveContext(actorId, 1, b))

		if !containsMsg(result.Messages, "混乱が解けた") {
			t.Errorf("expected confusion-cleared message, got %v", result.Messages)
		}
	})

	t.Run("継続中: 自傷（TargetSelf=true）と非自傷（TargetSelf=false）の両方が返る", func(t *testing.T) {
		gotSelfHit, gotNormal := false, false
		for range 300 {
			cf := rule.NewConfusion(3)
			st := status.NewStatusWith(nil, []status.OtherStatus{cf})
			actor := newPreMovePokemon("フシギダネ", st, 1)
			b, actorId := newPreMoveBattle(actor)

			result := phase.NewPreMovePhaseHandler().Handle(phase.NewPreMoveContext(actorId, 1, b))
			if result.NextPhase == phase.PhaseMoveResolve && result.ResolveSpec != nil {
				if result.ResolveSpec.TargetSelf {
					gotSelfHit = true
				} else {
					gotNormal = true
				}
			}
			if gotSelfHit && gotNormal {
				break
			}
		}
		if !gotSelfHit {
			t.Error("300試行で自傷（TargetSelf=true）が一度も発生しなかった")
		}
		if !gotNormal {
			t.Error("300試行で非自傷（TargetSelf=false）が一度も発生しなかった")
		}
	})

	t.Run("継続中 自傷時: 自傷メッセージあり", func(t *testing.T) {
		found := false
		for range 300 {
			cf := rule.NewConfusion(3)
			st := status.NewStatusWith(nil, []status.OtherStatus{cf})
			actor := newPreMovePokemon("フシギダネ", st, 1)
			b, actorId := newPreMoveBattle(actor)

			result := phase.NewPreMovePhaseHandler().Handle(phase.NewPreMoveContext(actorId, 1, b))
			if result.NextPhase == phase.PhaseMoveResolve && result.ResolveSpec != nil &&
				result.ResolveSpec.TargetSelf && containsMsg(result.Messages, "自分を攻撃した") {
				found = true
				break
			}
		}
		if !found {
			t.Error("300試行で自傷メッセージが確認できなかった")
		}
	})
}

// ── 8. まひ ───────────────────────────────────────────────────────────────────

func TestPreMoveHandle_Paralysis(t *testing.T) {
	t.Run("true と false（行動不能/行動可）の両方が返る", func(t *testing.T) {
		ms, err := status.NewMainStatus(status.Paralysis)
		if err != nil {
			t.Fatal(err)
		}
		gotEnd, gotDamage := false, false
		for range 300 {
			st := status.NewStatusWith(&ms, nil)
			actor := newPreMovePokemon("デンリュウ", st, 1)
			b, actorId := newPreMoveBattle(actor)

			result := phase.NewPreMovePhaseHandler().Handle(phase.NewPreMoveContext(actorId, 1, b))
			switch result.NextPhase {
			case phase.PhaseEnd:
				gotEnd = true
			case phase.PhaseMoveResolve:
				gotDamage = true
			}
			if gotEnd && gotDamage {
				break
			}
		}
		if !gotEnd {
			t.Error("300試行で麻痺行動不能（PhaseEnd）が一度も発生しなかった")
		}
		if !gotDamage {
			t.Error("300試行で行動成功（PhaseMoveResolve）が一度も発生しなかった")
		}
	})

	t.Run("行動不能時: 麻痺メッセージあり", func(t *testing.T) {
		ms, err := status.NewMainStatus(status.Paralysis)
		if err != nil {
			t.Fatal(err)
		}
		found := false
		for range 300 {
			st := status.NewStatusWith(&ms, nil)
			actor := newPreMovePokemon("デンリュウ", st, 1)
			b, actorId := newPreMoveBattle(actor)

			result := phase.NewPreMovePhaseHandler().Handle(phase.NewPreMoveContext(actorId, 1, b))
			if result.NextPhase == phase.PhaseEnd && containsMsg(result.Messages, "しびれて") {
				found = true
				break
			}
		}
		if !found {
			t.Error("300試行で麻痺メッセージが確認できなかった")
		}
	})
}
