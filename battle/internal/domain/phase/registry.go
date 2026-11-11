package phase

// Registry は各フェーズで実行するハンドラーを、発動タイミングごとに保持する。
// 特性・道具のハンドラーは ID をキーに登録し、フェーズ処理側で場のポケモンが持つ ID から引く。
// ダメージ補正・能力値補正は攻撃側／防御側で登録先を分け、登録した側として呼ばれた場合のみ実行される。
type Registry struct {
	entry handlerSet[EntryHandler] // 登場時（いかく等）
	exit  handlerSet[ExitHandler]  // 退場時（さいせいりょく等）

	resolveBaseHandlers []DamageModHandler              // 全ての攻撃に適用するダメージ補正（急所・タイプ相性・天気等）。登録順に実行
	resolve             sidedHandlers[DamageModHandler] // 特性・道具のダメージ補正（スナイパー・マルチスケイル等）
	resolveMoveHandlers map[int]MoveHandler             // 技ID ごとの技ハンドラー

	statModBaseHandlers []StatModHandler              // 全ての攻撃に適用する能力値補正（すなあらし等）。登録順に実行
	statMod             sidedHandlers[StatModHandler] // 特性・道具の能力値補正（ちからもち・とつげきチョッキ等）

	onHit    handlerSet[OnHitHandler]    // 被ダメージ時（じきゅうりょく・きのみ等）
	postMove handlerSet[PostMoveHandler] // ターン終了時（かそく等）
}

// NewRegistry は全ての登録先を空の状態で初期化した Registry を返す。
func NewRegistry() *Registry {
	r := &Registry{
		entry:               newHandlerSet[EntryHandler](),
		exit:                newHandlerSet[ExitHandler](),
		resolveBaseHandlers: []DamageModHandler{},
		resolve:             newSidedHandlers[DamageModHandler](),
		resolveMoveHandlers: map[int]MoveHandler{},
		statModBaseHandlers: []StatModHandler{},
		statMod:             newSidedHandlers[StatModHandler](),
		onHit:               newHandlerSet[OnHitHandler](),
		postMove:            newHandlerSet[PostMoveHandler](),
	}

	return r
}

// handlerSet は特性ID・道具ID ごとのハンドラーの組。
// 特性ID と道具ID はどちらも PokeAPI の ID で値が重なるため、map を分けて保持する。
type handlerSet[T any] struct {
	ability map[int]T
	item    map[int]T
}

func newHandlerSet[T any]() handlerSet[T] {
	return handlerSet[T]{ability: map[int]T{}, item: map[int]T{}}
}

// sidedHandlers は攻撃側・防御側それぞれの handlerSet。
// 同じ特性・道具でも、攻撃側で効くものと防御側で効くものを区別するために使う。
type sidedHandlers[T any] struct {
	attacker handlerSet[T]
	defender handlerSet[T]
}

func newSidedHandlers[T any]() sidedHandlers[T] {
	return sidedHandlers[T]{attacker: newHandlerSet[T](), defender: newHandlerSet[T]()}
}

// RegisterEntryAbility は登場時に発動する特性ハンドラーを登録する。同じIDで再登録した場合は上書きする。
func (r *Registry) RegisterEntryAbility(abilityId int, h EntryHandler) {
	r.entry.ability[abilityId] = h
}

// RegisterEntryItem は登場時に発動する道具ハンドラーを登録する。同じIDで再登録した場合は上書きする。
func (r *Registry) RegisterEntryItem(itemId int, h EntryHandler) {
	r.entry.item[itemId] = h
}

// RegisterExitAbility は退場時に発動する特性ハンドラーを登録する。同じIDで再登録した場合は上書きする。
func (r *Registry) RegisterExitAbility(abilityId int, h ExitHandler) {
	r.exit.ability[abilityId] = h
}

// RegisterExitItem は退場時に発動する道具ハンドラーを登録する。同じIDで再登録した場合は上書きする。
func (r *Registry) RegisterExitItem(itemId int, h ExitHandler) {
	r.exit.item[itemId] = h
}

// AddResolveBase は全ての攻撃に適用するダメージ補正ハンドラーを末尾に追加する。登録順に実行される。
// CritHandler は先頭に登録すること（後続ハンドラーが ctx.IsCrit を参照するため）。
func (r *Registry) AddResolveBase(h DamageModHandler) {
	r.resolveBaseHandlers = append(r.resolveBaseHandlers, h)
}

// RegisterResolveAttackerAbility は攻撃側の特性として発動するダメージ補正ハンドラーを登録する。
func (r *Registry) RegisterResolveAttackerAbility(abilityId int, h DamageModHandler) {
	r.resolve.attacker.ability[abilityId] = h
}

// RegisterResolveAttackerItem は攻撃側の道具として発動するダメージ補正ハンドラーを登録する。
func (r *Registry) RegisterResolveAttackerItem(itemId int, h DamageModHandler) {
	r.resolve.attacker.item[itemId] = h
}

// RegisterResolveDefenderAbility は防御側の特性として発動するダメージ補正ハンドラーを登録する。
func (r *Registry) RegisterResolveDefenderAbility(abilityId int, h DamageModHandler) {
	r.resolve.defender.ability[abilityId] = h
}

// RegisterResolveDefenderItem は防御側の道具として発動するダメージ補正ハンドラーを登録する。
func (r *Registry) RegisterResolveDefenderItem(itemId int, h DamageModHandler) {
	r.resolve.defender.item[itemId] = h
}

// RegisterResolveMove は技IDに対応する技ハンドラーを登録する。同じIDで再登録した場合は上書きする。
func (r *Registry) RegisterResolveMove(moveId int, h MoveHandler) {
	r.resolveMoveHandlers[moveId] = h
}

// AddStatModBase は全ての攻撃に適用する能力値補正ハンドラーを末尾に追加する。登録順に実行される。
func (r *Registry) AddStatModBase(h StatModHandler) {
	r.statModBaseHandlers = append(r.statModBaseHandlers, h)
}

// RegisterStatModAttackerAbility は攻撃側の特性として発動する能力値補正ハンドラーを登録する。
func (r *Registry) RegisterStatModAttackerAbility(abilityId int, h StatModHandler) {
	r.statMod.attacker.ability[abilityId] = h
}

// RegisterStatModAttackerItem は攻撃側の道具として発動する能力値補正ハンドラーを登録する。
func (r *Registry) RegisterStatModAttackerItem(itemId int, h StatModHandler) {
	r.statMod.attacker.item[itemId] = h
}

// RegisterStatModDefenderAbility は防御側の特性として発動する能力値補正ハンドラーを登録する。
func (r *Registry) RegisterStatModDefenderAbility(abilityId int, h StatModHandler) {
	r.statMod.defender.ability[abilityId] = h
}

// RegisterStatModDefenderItem は防御側の道具として発動する能力値補正ハンドラーを登録する。
func (r *Registry) RegisterStatModDefenderItem(itemId int, h StatModHandler) {
	r.statMod.defender.item[itemId] = h
}

// RegisterOnHitAbility は被ダメージ時に発動する特性ハンドラーを登録する。同じIDで再登録した場合は上書きする。
func (r *Registry) RegisterOnHitAbility(abilityId int, h OnHitHandler) {
	r.onHit.ability[abilityId] = h
}

// RegisterOnHitItem は被ダメージ時に発動する道具ハンドラーを登録する。同じIDで再登録した場合は上書きする。
func (r *Registry) RegisterOnHitItem(itemId int, h OnHitHandler) {
	r.onHit.item[itemId] = h
}

// RegisterPostMoveAbility はターン終了時に発動する特性ハンドラーを登録する。同じIDで再登録した場合は上書きする。
func (r *Registry) RegisterPostMoveAbility(abilityId int, h PostMoveHandler) {
	r.postMove.ability[abilityId] = h
}

// RegisterPostMoveItem はターン終了時に発動する道具ハンドラーを登録する。同じIDで再登録した場合は上書きする。
func (r *Registry) RegisterPostMoveItem(itemId int, h PostMoveHandler) {
	r.postMove.item[itemId] = h
}
