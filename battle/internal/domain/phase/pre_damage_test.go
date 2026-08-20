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

func newPreDamagePokemon(name string, st status.Status, moveIds ...int) *pokemon.Pokemon {
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

// newPreDamageBattle は actor を p1 のアクティブポケモンとしてバトルを構築する。
// 返値の actorId は PreDamageContext に渡す p1 のプレイヤー ID。
func newPreDamageBattle(actor *pokemon.Pokemon) (*battle.Battle, string) {
	dummy1 := newPreDamagePokemon("ダミー1", status.NewStatus(), 1)
	dummy2 := newPreDamagePokemon("ダミー2", status.NewStatus(), 1)
	opp1 := newPreDamagePokemon("あいて1", status.NewStatus(), 1)
	opp2 := newPreDamagePokemon("あいて2", status.NewStatus(), 1)
	opp3 := newPreDamagePokemon("あいて3", status.NewStatus(), 1)

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

func TestPreDamageHandle_NoStatus(t *testing.T) {
	t.Run("状態異常なし: PhaseDamage に進む", func(t *testing.T) {
		actor := newPreDamagePokemon("ピカチュウ", status.NewStatus(), 1)
		b, actorId := newPreDamageBattle(actor)

		result := phase.NewPreDamagePhaseHandler().Handle(phase.NewPreDamageContext(actorId, 1, b))

		if result.NextPhase != phase.PhaseDamage {
			t.Errorf("expected PhaseDamage, got %v", result.NextPhase)
		}
	})

	t.Run("状態異常なし: メッセージなし", func(t *testing.T) {
		actor := newPreDamagePokemon("ピカチュウ", status.NewStatus(), 1)
		b, actorId := newPreDamageBattle(actor)

		result := phase.NewPreDamagePhaseHandler().Handle(phase.NewPreDamageContext(actorId, 1, b))

		if len(result.Messages) != 0 {
			t.Errorf("expected no messages, got %v", result.Messages)
		}
	})
}

// ── 2. ねむり ─────────────────────────────────────────────────────────────────

func TestPreDamageHandle_Sleep(t *testing.T) {
	t.Run("眠り継続中: PhaseEnd", func(t *testing.T) {
		ms, err := status.NewSleep(vo.NewCount(2))
		if err != nil {
			t.Fatal(err)
		}
		st := status.NewStatusWith(&ms, nil)
		actor := newPreDamagePokemon("カビゴン", st, 1)
		b, actorId := newPreDamageBattle(actor)

		result := phase.NewPreDamagePhaseHandler().Handle(phase.NewPreDamageContext(actorId, 1, b))

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
		actor := newPreDamagePokemon("カビゴン", st, 1)
		b, actorId := newPreDamageBattle(actor)

		result := phase.NewPreDamagePhaseHandler().Handle(phase.NewPreDamageContext(actorId, 1, b))

		if !containsMsg(result.Messages, "ぐうぐう眠っている") {
			t.Errorf("expected sleep message, got %v", result.Messages)
		}
	})
}

// ── 3. こおり ─────────────────────────────────────────────────────────────────

// NewMainStatus(Freeze) は count=0 で生成されるため IsFrozen が即時 false を返し、解凍扱いとなる。
func TestPreDamageHandle_Freeze(t *testing.T) {
	t.Run("解凍: PhaseDamage に進む", func(t *testing.T) {
		ms, err := status.NewMainStatus(status.Freeze)
		if err != nil {
			t.Fatal(err)
		}
		st := status.NewStatusWith(&ms, nil)
		actor := newPreDamagePokemon("ラプラス", st, 1)
		b, actorId := newPreDamageBattle(actor)

		result := phase.NewPreDamagePhaseHandler().Handle(phase.NewPreDamageContext(actorId, 1, b))

		if result.NextPhase != phase.PhaseDamage {
			t.Errorf("expected PhaseDamage, got %v", result.NextPhase)
		}
	})

	t.Run("解凍: 解凍メッセージあり", func(t *testing.T) {
		ms, err := status.NewMainStatus(status.Freeze)
		if err != nil {
			t.Fatal(err)
		}
		st := status.NewStatusWith(&ms, nil)
		actor := newPreDamagePokemon("ラプラス", st, 1)
		b, actorId := newPreDamageBattle(actor)

		result := phase.NewPreDamagePhaseHandler().Handle(phase.NewPreDamageContext(actorId, 1, b))

		if !containsMsg(result.Messages, "こおりが溶けた") {
			t.Errorf("expected thaw message, got %v", result.Messages)
		}
	})
}

// ── 4. ひるみ ─────────────────────────────────────────────────────────────────

func TestPreDamageHandle_Flinch(t *testing.T) {
	t.Run("ひるみあり: PhaseEnd", func(t *testing.T) {
		fl := rule.NewFlinch()
		st := status.NewStatusWith(nil, []status.OtherStatus{fl})
		actor := newPreDamagePokemon("ピカチュウ", st, 1)
		b, actorId := newPreDamageBattle(actor)

		result := phase.NewPreDamagePhaseHandler().Handle(phase.NewPreDamageContext(actorId, 1, b))

		if result.NextPhase != phase.PhaseEnd {
			t.Errorf("expected PhaseEnd, got %v", result.NextPhase)
		}
	})

	t.Run("ひるみあり: ひるみメッセージあり", func(t *testing.T) {
		fl := rule.NewFlinch()
		st := status.NewStatusWith(nil, []status.OtherStatus{fl})
		actor := newPreDamagePokemon("ピカチュウ", st, 1)
		b, actorId := newPreDamageBattle(actor)

		result := phase.NewPreDamagePhaseHandler().Handle(phase.NewPreDamageContext(actorId, 1, b))

		if !containsMsg(result.Messages, "ひるんで") {
			t.Errorf("expected flinch message, got %v", result.Messages)
		}
	})
}

// ── 5. アンコール ─────────────────────────────────────────────────────────────

func TestPreDamageHandle_Encore(t *testing.T) {
	t.Run("アンコール技以外を選択: PhaseEnd", func(t *testing.T) {
		enc := rule.NewEncore(1)
		st := status.NewStatusWith(nil, []status.OtherStatus{enc})
		actor := newPreDamagePokemon("ピクシー", st, 1, 2)
		b, actorId := newPreDamageBattle(actor)

		result := phase.NewPreDamagePhaseHandler().Handle(phase.NewPreDamageContext(actorId, 2, b))

		if result.NextPhase != phase.PhaseEnd {
			t.Errorf("expected PhaseEnd, got %v", result.NextPhase)
		}
	})

	t.Run("アンコール技以外を選択: アンコールメッセージあり", func(t *testing.T) {
		enc := rule.NewEncore(1)
		st := status.NewStatusWith(nil, []status.OtherStatus{enc})
		actor := newPreDamagePokemon("ピクシー", st, 1, 2)
		b, actorId := newPreDamageBattle(actor)

		result := phase.NewPreDamagePhaseHandler().Handle(phase.NewPreDamageContext(actorId, 2, b))

		if !containsMsg(result.Messages, "アンコール") {
			t.Errorf("expected encore message, got %v", result.Messages)
		}
	})

	t.Run("アンコール技を選択: PhaseDamage に進む", func(t *testing.T) {
		enc := rule.NewEncore(1)
		st := status.NewStatusWith(nil, []status.OtherStatus{enc})
		actor := newPreDamagePokemon("ピクシー", st, 1, 2)
		b, actorId := newPreDamageBattle(actor)

		result := phase.NewPreDamagePhaseHandler().Handle(phase.NewPreDamageContext(actorId, 1, b))

		if result.NextPhase != phase.PhaseDamage {
			t.Errorf("expected PhaseDamage, got %v", result.NextPhase)
		}
	})

	t.Run("アンコール技を選択: メッセージなし", func(t *testing.T) {
		enc := rule.NewEncore(1)
		st := status.NewStatusWith(nil, []status.OtherStatus{enc})
		actor := newPreDamagePokemon("ピクシー", st, 1, 2)
		b, actorId := newPreDamageBattle(actor)

		result := phase.NewPreDamagePhaseHandler().Handle(phase.NewPreDamageContext(actorId, 1, b))

		if len(result.Messages) != 0 {
			t.Errorf("expected no messages, got %v", result.Messages)
		}
	})
}

// ── 6. かなしばり ─────────────────────────────────────────────────────────────

func TestPreDamageHandle_MoveDisabled(t *testing.T) {
	t.Run("封じられた技を選択: PhaseEnd", func(t *testing.T) {
		md := rule.NewMoveDisabled(4, 1)
		st := status.NewStatusWith(nil, []status.OtherStatus{md})
		actor := newPreDamagePokemon("ゲンガー", st, 1, 2)
		b, actorId := newPreDamageBattle(actor)

		result := phase.NewPreDamagePhaseHandler().Handle(phase.NewPreDamageContext(actorId, 1, b))

		if result.NextPhase != phase.PhaseEnd {
			t.Errorf("expected PhaseEnd, got %v", result.NextPhase)
		}
	})

	t.Run("封じられた技を選択: かなしばりメッセージあり", func(t *testing.T) {
		md := rule.NewMoveDisabled(4, 1)
		st := status.NewStatusWith(nil, []status.OtherStatus{md})
		actor := newPreDamagePokemon("ゲンガー", st, 1, 2)
		b, actorId := newPreDamageBattle(actor)

		result := phase.NewPreDamagePhaseHandler().Handle(phase.NewPreDamageContext(actorId, 1, b))

		if !containsMsg(result.Messages, "かなしばり") {
			t.Errorf("expected move-disabled message, got %v", result.Messages)
		}
	})

	t.Run("封じられていない技を選択: PhaseDamage に進む", func(t *testing.T) {
		md := rule.NewMoveDisabled(4, 1)
		st := status.NewStatusWith(nil, []status.OtherStatus{md})
		actor := newPreDamagePokemon("ゲンガー", st, 1, 2)
		b, actorId := newPreDamageBattle(actor)

		result := phase.NewPreDamagePhaseHandler().Handle(phase.NewPreDamageContext(actorId, 2, b))

		if result.NextPhase != phase.PhaseDamage {
			t.Errorf("expected PhaseDamage, got %v", result.NextPhase)
		}
	})
}

// ── 7. こんらん ───────────────────────────────────────────────────────────────

func TestPreDamageHandle_Confusion(t *testing.T) {
	t.Run("残り1ターン: 今ターンで解除 → PhaseDamage", func(t *testing.T) {
		cf := rule.NewConfusion(1)
		st := status.NewStatusWith(nil, []status.OtherStatus{cf})
		actor := newPreDamagePokemon("フシギダネ", st, 1)
		b, actorId := newPreDamageBattle(actor)

		result := phase.NewPreDamagePhaseHandler().Handle(phase.NewPreDamageContext(actorId, 1, b))

		if result.NextPhase != phase.PhaseDamage {
			t.Errorf("expected PhaseDamage, got %v", result.NextPhase)
		}
	})

	t.Run("残り1ターン: 解除メッセージあり", func(t *testing.T) {
		cf := rule.NewConfusion(1)
		st := status.NewStatusWith(nil, []status.OtherStatus{cf})
		actor := newPreDamagePokemon("フシギダネ", st, 1)
		b, actorId := newPreDamageBattle(actor)

		result := phase.NewPreDamagePhaseHandler().Handle(phase.NewPreDamageContext(actorId, 1, b))

		if !containsMsg(result.Messages, "混乱が解けた") {
			t.Errorf("expected confusion-cleared message, got %v", result.Messages)
		}
	})

	t.Run("継続中: true と false（自傷/非自傷）の両方が返る", func(t *testing.T) {
		gotEnd, gotDamage := false, false
		for range 300 {
			cf := rule.NewConfusion(3)
			st := status.NewStatusWith(nil, []status.OtherStatus{cf})
			actor := newPreDamagePokemon("フシギダネ", st, 1)
			b, actorId := newPreDamageBattle(actor)

			result := phase.NewPreDamagePhaseHandler().Handle(phase.NewPreDamageContext(actorId, 1, b))
			switch result.NextPhase {
			case phase.PhaseEnd:
				gotEnd = true
			case phase.PhaseDamage:
				gotDamage = true
			}
			if gotEnd && gotDamage {
				break
			}
		}
		if !gotEnd {
			t.Error("300試行で自傷（PhaseEnd）が一度も発生しなかった")
		}
		if !gotDamage {
			t.Error("300試行で非自傷（PhaseDamage）が一度も発生しなかった")
		}
	})

	t.Run("継続中 自傷時: 自傷メッセージあり", func(t *testing.T) {
		found := false
		for range 300 {
			cf := rule.NewConfusion(3)
			st := status.NewStatusWith(nil, []status.OtherStatus{cf})
			actor := newPreDamagePokemon("フシギダネ", st, 1)
			b, actorId := newPreDamageBattle(actor)

			result := phase.NewPreDamagePhaseHandler().Handle(phase.NewPreDamageContext(actorId, 1, b))
			if result.NextPhase == phase.PhaseEnd && containsMsg(result.Messages, "自分を攻撃した") {
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

func TestPreDamageHandle_Paralysis(t *testing.T) {
	t.Run("true と false（行動不能/行動可）の両方が返る", func(t *testing.T) {
		ms, err := status.NewMainStatus(status.Paralysis)
		if err != nil {
			t.Fatal(err)
		}
		gotEnd, gotDamage := false, false
		for range 300 {
			st := status.NewStatusWith(&ms, nil)
			actor := newPreDamagePokemon("デンリュウ", st, 1)
			b, actorId := newPreDamageBattle(actor)

			result := phase.NewPreDamagePhaseHandler().Handle(phase.NewPreDamageContext(actorId, 1, b))
			switch result.NextPhase {
			case phase.PhaseEnd:
				gotEnd = true
			case phase.PhaseDamage:
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
			t.Error("300試行で行動成功（PhaseDamage）が一度も発生しなかった")
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
			actor := newPreDamagePokemon("デンリュウ", st, 1)
			b, actorId := newPreDamageBattle(actor)

			result := phase.NewPreDamagePhaseHandler().Handle(phase.NewPreDamageContext(actorId, 1, b))
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
