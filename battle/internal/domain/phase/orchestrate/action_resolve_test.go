package orchestrate_test

import (
	"strings"
	"testing"

	"pob/battle/internal/domain/ability"
	"pob/battle/internal/domain/battle"
	"pob/battle/internal/domain/hp"
	"pob/battle/internal/domain/move"
	"pob/battle/internal/domain/nature"
	"pob/battle/internal/domain/phase"
	"pob/battle/internal/domain/phase/orchestrate"
	"pob/battle/internal/domain/player"
	"pob/battle/internal/domain/pokemon"
	"pob/battle/internal/domain/pp"
	"pob/battle/internal/domain/ptype"
	"pob/battle/internal/domain/rank"
	"pob/battle/internal/domain/vo"
)

// ── helpers ───────────────────────────────────────────────────────────────────

// newTestPokemon は最低限の値を持つポケモンを生成するヘルパー。
// abilityId=0 / moveId=1(PP10) を固定で持ち、HP が満タンなので気絶していない。
func newTestPokemon(name string, speed int) *pokemon.Pokemon {
	a := ability.NewAbility(0, "")
	m := move.NewMove(1, pp.NewPP(vo.NewCount(10), 10), 0, 0, 0, move.DamageClassStatus, ptype.Normal)
	return pokemon.NewPokemon(
		1, name, 0, "",
		[2]ptype.Type{ptype.Normal, ptype.Normal},
		pokemon.BaseStats{},
		pokemon.RealStats{Speed: speed},
		nature.Nature{},
		&a,
		[4]*move.Move{&m},
		hp.NewHP(100),
		rank.NewRank(),
		nil, nil, nil, nil,
		false,
	)
}

// newTestPlayer は party[0..2] を選出済みの状態で Player を生成するヘルパー。
func newTestPlayer(id, name string, poke0, poke1, poke2 *pokemon.Pokemon) *player.Player {
	party := [6]*pokemon.Pokemon{poke0, poke1, poke2}
	pl := player.NewPlayer(id, name, party, nil)
	if err := pl.Select([3]int{0, 1, 2}); err != nil {
		panic(err)
	}
	return pl
}

// newTestOrchestrator は空の Registry を使った Orchestrator を生成するヘルパー。
func newTestOrchestrator() *orchestrate.ActionResolveOrchestrator {
	r := phase.NewRegistry()
	return orchestrate.NewActionResolveOrchestrator(
		phase.NewExitPhaseHandler(r),
		phase.NewEntryPhaseHandler(r),
		&phase.MoveSelectPhaseHandler{},
		&phase.ForfeitPhaseHandler{},
	)
}

// containsSubstr はメッセージスライス内にサブストリングを含む要素があるか調べる。
func containsSubstr(msgs []string, sub string) bool {
	for _, m := range msgs {
		if strings.Contains(m, sub) {
			return true
		}
	}
	return false
}

// indexOfContaining はサブストリングを含む最初の要素のインデックスを返す（なければ -1）。
func indexOfContaining(msgs []string, sub string) int {
	for i, m := range msgs {
		if strings.Contains(m, sub) {
			return i
		}
	}
	return -1
}

// ── 1. バリデーション ─────────────────────────────────────────────────────────

func TestHandle_Validation_NoPendingAction(t *testing.T) {
	p1 := newTestPlayer("p1", "サトシ", newTestPokemon("アーボ", 100), newTestPokemon("バタフリー", 80), newTestPokemon("コダック", 60))
	p2 := newTestPlayer("p2", "シゲル", newTestPokemon("ドードー", 90), newTestPokemon("エビワラー", 70), newTestPokemon("フーディン", 50))

	b := battle.NewBattle(p1, p2)
	orc := newTestOrchestrator()

	_, err := orc.Handle(b)
	if err == nil {
		t.Fatal("保留行動なし: エラーを期待したが nil")
	}
}

// p1 が Forfeit() 後に Switch(1) を呼ぶと、Switch が pendingForfeit を消去しないため
// forfeit と switch の 2 つが同時に pending になる。
func TestHandle_Validation_MultiplePendingActions(t *testing.T) {
	p1 := newTestPlayer("p1", "サトシ", newTestPokemon("アーボ", 100), newTestPokemon("バタフリー", 80), newTestPokemon("コダック", 60))
	p2 := newTestPlayer("p2", "シゲル", newTestPokemon("ドードー", 90), newTestPokemon("エビワラー", 70), newTestPokemon("フーディン", 50))

	p1.Forfeit()
	if err := p1.Switch(1); err != nil {
		t.Fatalf("Switch failed: %v", err)
	}
	if err := p2.Switch(1); err != nil {
		t.Fatalf("Switch failed: %v", err)
	}

	b := battle.NewBattle(p1, p2)
	orc := newTestOrchestrator()

	_, err := orc.Handle(b)
	if err == nil {
		t.Fatal("複数保留行動: エラーを期待したが nil")
	}
}

// ── 2. サレンダー ─────────────────────────────────────────────────────────────

func TestHandle_Forfeit_Single(t *testing.T) {
	p1 := newTestPlayer("p1", "サトシ", newTestPokemon("アーボ", 100), newTestPokemon("バタフリー", 80), newTestPokemon("コダック", 60))
	p2 := newTestPlayer("p2", "シゲル", newTestPokemon("ドードー", 90), newTestPokemon("エビワラー", 70), newTestPokemon("フーディン", 50))

	p1.Forfeit()
	if err := p2.Switch(1); err != nil {
		t.Fatalf("Switch failed: %v", err)
	}

	b := battle.NewBattle(p1, p2)
	orc := newTestOrchestrator()

	result, err := orc.Handle(b)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// 勝者が p2 であること
	if b.Winner() != p2 {
		t.Error("Winner は p2 のはず")
	}

	// p1（降参した側）に "降参" メッセージが入ること
	if !containsSubstr(result["p1"], "降参しました") {
		t.Errorf("p1 に降参メッセージなし: %v", result["p1"])
	}

	// p2（勝利した側）に "勝利" メッセージが入ること
	if !containsSubstr(result["p2"], "勝利です") {
		t.Errorf("p2 に勝利メッセージなし: %v", result["p2"])
	}
}

func TestHandle_Forfeit_Both_Draw(t *testing.T) {
	p1 := newTestPlayer("p1", "サトシ", newTestPokemon("アーボ", 100), newTestPokemon("バタフリー", 80), newTestPokemon("コダック", 60))
	p2 := newTestPlayer("p2", "シゲル", newTestPokemon("ドードー", 90), newTestPokemon("エビワラー", 70), newTestPokemon("フーディン", 50))

	p1.Forfeit()
	p2.Forfeit()

	b := battle.NewBattle(p1, p2)
	orc := newTestOrchestrator()

	result, err := orc.Handle(b)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !b.IsDraw() {
		t.Error("両者サレンダーは引き分けのはず")
	}
	if !containsSubstr(result["p1"], "引き分け") {
		t.Errorf("p1 に引き分けメッセージなし: %v", result["p1"])
	}
	if !containsSubstr(result["p2"], "引き分け") {
		t.Errorf("p2 に引き分けメッセージなし: %v", result["p2"])
	}
}

// サレンダー時は交代処理ブロックに到達する前に早期 return されるため、
// p2 の交代は実行されず activeSlot が変わらないことを確認する。
func TestHandle_Forfeit_SkipsSwitchProcessing(t *testing.T) {
	p1 := newTestPlayer("p1", "サトシ", newTestPokemon("アーボ", 100), newTestPokemon("バタフリー", 80), newTestPokemon("コダック", 60))
	p2 := newTestPlayer("p2", "シゲル", newTestPokemon("ドードー", 90), newTestPokemon("エビワラー", 70), newTestPokemon("フーディン", 50))

	p2BeforeSwitch := p2.Active()

	p1.Forfeit()
	if err := p2.Switch(1); err != nil {
		t.Fatalf("Switch failed: %v", err)
	}

	b := battle.NewBattle(p1, p2)
	orc := newTestOrchestrator()

	if _, err := orc.Handle(b); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if p2.Active() != p2BeforeSwitch {
		t.Error("サレンダー早期 return により p2 の交代は実行されないはず")
	}
}

// ── 3. 交代（単独） ───────────────────────────────────────────────────────────

func TestHandle_Switch_Single_Messages(t *testing.T) {
	outgoing := newTestPokemon("フシギダネ", 100)
	incoming := newTestPokemon("ヒトカゲ", 90)
	p1 := newTestPlayer("p1", "サトシ", outgoing, incoming, newTestPokemon("ゼニガメ", 80))
	p2 := newTestPlayer("p2", "シゲル", newTestPokemon("ピカチュウ", 70), newTestPokemon("イーブイ", 60), newTestPokemon("カビゴン", 50))

	if err := p1.Switch(1); err != nil {
		t.Fatalf("Switch failed: %v", err)
	}
	if err := p2.SelectMove(1); err != nil {
		t.Fatalf("SelectMove failed: %v", err)
	}

	b := battle.NewBattle(p1, p2)
	orc := newTestOrchestrator()

	result, err := orc.Handle(b)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// p1（自分視点）: "戻れ！" と "行け！" が入ること
	if !containsSubstr(result["p1"], "戻れ！フシギダネ！") {
		t.Errorf("p1 に戻れメッセージなし: %v", result["p1"])
	}
	if !containsSubstr(result["p1"], "行け！ヒトカゲ！") {
		t.Errorf("p1 に行けメッセージなし: %v", result["p1"])
	}

	// p2（相手視点）: "引っ込めた" と "繰り出した" が入ること
	if !containsSubstr(result["p2"], "相手はフシギダネを引っ込めた！") {
		t.Errorf("p2 に引っ込めたメッセージなし: %v", result["p2"])
	}
	if !containsSubstr(result["p2"], "相手はヒトカゲを繰り出した！") {
		t.Errorf("p2 に繰り出したメッセージなし: %v", result["p2"])
	}
}

// commitSwitch によって p1 の activeSlot が正しく更新されること。
func TestHandle_Switch_Single_ActiveSlot(t *testing.T) {
	incoming := newTestPokemon("ヒトカゲ", 90)
	p1 := newTestPlayer("p1", "サトシ", newTestPokemon("フシギダネ", 100), incoming, newTestPokemon("ゼニガメ", 80))
	p2 := newTestPlayer("p2", "シゲル", newTestPokemon("ピカチュウ", 70), newTestPokemon("イーブイ", 60), newTestPokemon("カビゴン", 50))

	if err := p1.Switch(1); err != nil {
		t.Fatalf("Switch failed: %v", err)
	}
	if err := p2.SelectMove(1); err != nil {
		t.Fatalf("SelectMove failed: %v", err)
	}

	b := battle.NewBattle(p1, p2)
	orc := newTestOrchestrator()

	if _, err := orc.Handle(b); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if p1.Active() != incoming {
		t.Error("commitSwitch 後は incoming ポケモンが active になるはず")
	}
}

// ── 4. 交代（同時・速度順） ───────────────────────────────────────────────────

// p1 の退場ポケモン（speed=100）が p2（speed=50）より速い場合、
// p1 の退場メッセージが p2 の退場メッセージより先に result に入ること。
func TestHandle_Switch_Both_SpeedOrder(t *testing.T) {
	fastPoke := newTestPokemon("はやいポケモン", 100)
	fastIn := newTestPokemon("はやいひかえ", 90)
	slowPoke := newTestPokemon("おそいポケモン", 50)
	slowIn := newTestPokemon("おそいひかえ", 40)

	p1 := newTestPlayer("p1", "サトシ", fastPoke, fastIn, newTestPokemon("コイキング", 10))
	p2 := newTestPlayer("p2", "シゲル", slowPoke, slowIn, newTestPokemon("ヤドン", 10))

	if err := p1.Switch(1); err != nil {
		t.Fatalf("p1.Switch: %v", err)
	}
	if err := p2.Switch(1); err != nil {
		t.Fatalf("p2.Switch: %v", err)
	}

	b := battle.NewBattle(p1, p2)
	orc := newTestOrchestrator()

	result, err := orc.Handle(b)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// result["p1"] において はやいポケモンの退場が おそいポケモンの退場より先に現れること
	msgs := result["p1"]
	fastIdx := indexOfContaining(msgs, "戻れ！はやいポケモン！")
	slowIdx := indexOfContaining(msgs, "相手はおそいポケモンを引っ込めた！")

	if fastIdx < 0 {
		t.Fatalf("はやいポケモン退場メッセージなし: %v", msgs)
	}
	if slowIdx < 0 {
		t.Fatalf("おそいポケモン退場メッセージなし: %v", msgs)
	}
	if fastIdx >= slowIdx {
		t.Errorf("速い方（はやいポケモン）の処理が先のはず: fastIdx=%d, slowIdx=%d, msgs=%v", fastIdx, slowIdx, msgs)
	}
}

// 速度が同じ場合は stable sort により入力順（p1 → p2）が保たれること。
func TestHandle_Switch_Both_SameSpeed_InputOrder(t *testing.T) {
	p1Poke := newTestPokemon("P1先発", 100)
	p1In := newTestPokemon("P1ひかえ", 90)
	p2Poke := newTestPokemon("P2先発", 100) // 同速
	p2In := newTestPokemon("P2ひかえ", 90)

	p1 := newTestPlayer("p1", "サトシ", p1Poke, p1In, newTestPokemon("コイキング", 10))
	p2 := newTestPlayer("p2", "シゲル", p2Poke, p2In, newTestPokemon("ヤドン", 10))

	if err := p1.Switch(1); err != nil {
		t.Fatalf("p1.Switch: %v", err)
	}
	if err := p2.Switch(1); err != nil {
		t.Fatalf("p2.Switch: %v", err)
	}

	b := battle.NewBattle(p1, p2)
	orc := newTestOrchestrator()

	result, err := orc.Handle(b)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// result["p1"] において p1 自身の退場（入力順）が p2 の退場（opponent 視点）より先に現れること
	msgs := result["p1"]
	p1SelfIdx := indexOfContaining(msgs, "戻れ！P1先発！")
	p2OppIdx := indexOfContaining(msgs, "相手はP2先発を引っ込めた！")

	if p1SelfIdx < 0 || p2OppIdx < 0 {
		t.Fatalf("期待するメッセージなし: %v", msgs)
	}
	if p1SelfIdx >= p2OppIdx {
		t.Errorf("同速では入力順（p1 先）が保たれるはず: p1SelfIdx=%d, p2OppIdx=%d", p1SelfIdx, p2OppIdx)
	}
}

// ── 5. 交代 + 技の組み合わせ ─────────────────────────────────────────────────

// 片方が交代・片方が技を選択した場合、交代メッセージだけ result に入り、
// 技は b.PendingMoves() に積まれること（このフェーズでは実行されない）。
func TestHandle_SwitchAndMove(t *testing.T) {
	outgoing := newTestPokemon("フシギダネ", 100)
	incoming := newTestPokemon("ヒトカゲ", 90)
	p1 := newTestPlayer("p1", "サトシ", outgoing, incoming, newTestPokemon("ゼニガメ", 80))
	p2 := newTestPlayer("p2", "シゲル", newTestPokemon("ピカチュウ", 70), newTestPokemon("イーブイ", 60), newTestPokemon("カビゴン", 50))

	if err := p1.Switch(1); err != nil {
		t.Fatalf("Switch failed: %v", err)
	}
	if err := p2.SelectMove(1); err != nil {
		t.Fatalf("SelectMove failed: %v", err)
	}

	b := battle.NewBattle(p1, p2)
	orc := newTestOrchestrator()

	result, err := orc.Handle(b)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// 交代メッセージが result に含まれること
	if !containsSubstr(result["p1"], "戻れ！フシギダネ！") {
		t.Errorf("交代メッセージなし: %v", result["p1"])
	}

	// 技は PendingMoves にキューされること
	pending := b.PendingMoves()
	if len(pending) != 1 {
		t.Fatalf("PendingMoves は 1 件のはず: got %d", len(pending))
	}
	if pending[0].MoveId != 1 {
		t.Errorf("PendingMoves[0].MoveId = %d, want 1", pending[0].MoveId)
	}

	// 技の result メッセージはこのフェーズでは生成されない
	// → result["p2"] は交代に対する相手視点メッセージのみ（2件）
	if len(result["p2"]) != 2 {
		t.Errorf("result[p2] は 2 件のはず（交代の相手視点のみ）: %v", result["p2"])
	}
}

// ── 6. メッセージのマージ ─────────────────────────────────────────────────────

// exitHandler / entryHandler が返す map が mergeMessages で result に
// 正しい順序で統合されること。
func TestHandle_Messages_Merge_Order(t *testing.T) {
	outgoing := newTestPokemon("でていくポケモン", 100)
	incoming := newTestPokemon("はいってくるポケモン", 80)
	p1 := newTestPlayer("p1", "サトシ", outgoing, incoming, newTestPokemon("コイキング", 10))
	p2 := newTestPlayer("p2", "シゲル", newTestPokemon("あいてポケモン", 70), newTestPokemon("ヤドン", 60), newTestPokemon("ズバット", 50))

	if err := p1.Switch(1); err != nil {
		t.Fatalf("Switch failed: %v", err)
	}
	if err := p2.SelectMove(1); err != nil {
		t.Fatalf("SelectMove failed: %v", err)
	}

	b := battle.NewBattle(p1, p2)
	orc := newTestOrchestrator()

	result, err := orc.Handle(b)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// p1（自分視点）: 退場メッセージ → 入場メッセージの順
	p1Msgs := result["p1"]
	if len(p1Msgs) < 2 {
		t.Fatalf("result[p1] は 2 件以上のはず: %v", p1Msgs)
	}
	if p1Msgs[0] != "戻れ！でていくポケモン！" {
		t.Errorf("result[p1][0] = %q, want \"戻れ！でていくポケモン！\"", p1Msgs[0])
	}
	if p1Msgs[1] != "行け！はいってくるポケモン！" {
		t.Errorf("result[p1][1] = %q, want \"行け！はいってくるポケモン！\"", p1Msgs[1])
	}

	// p2（相手視点）: 引っ込めた → 繰り出したの順
	p2Msgs := result["p2"]
	if len(p2Msgs) < 2 {
		t.Fatalf("result[p2] は 2 件以上のはず: %v", p2Msgs)
	}
	if p2Msgs[0] != "相手はでていくポケモンを引っ込めた！" {
		t.Errorf("result[p2][0] = %q, want \"相手はでていくポケモンを引っ込めた！\"", p2Msgs[0])
	}
	if p2Msgs[1] != "相手ははいってくるポケモンを繰り出した！" {
		t.Errorf("result[p2][1] = %q, want \"相手ははいってくるポケモンを繰り出した！\"", p2Msgs[1])
	}
}
