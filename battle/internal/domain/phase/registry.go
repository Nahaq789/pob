package phase

type Registry struct {
	entryAbilityHandlers   map[int]EntryHandler
	entryItemHandlers      map[int]EntryHandler
	exitAbilityHandlers    map[int]ExitHandler
	exitItemHandlers       map[int]ExitHandler
	resolveBaseHandlers    []DamageModHandler
	resolveAbilityHandlers map[int]DamageModHandler
	resolveItemHandlers    map[int]DamageModHandler
	resolveMoveHandlers    map[int]MoveHandler

	onHitAbilityHandlers map[int]OnHitHandler
	onHitItemHandlers    map[int]OnHitHandler

	postMoveAbilityHandlers map[int]PostMoveHandler
	postMoveItemHandlers    map[int]PostMoveHandler
}

func NewRegistry() *Registry {
	r := &Registry{
		entryAbilityHandlers:    map[int]EntryHandler{},
		entryItemHandlers:       map[int]EntryHandler{},
		exitAbilityHandlers:     map[int]ExitHandler{},
		exitItemHandlers:        map[int]ExitHandler{},
		resolveBaseHandlers:     []DamageModHandler{},
		resolveAbilityHandlers:  map[int]DamageModHandler{},
		resolveItemHandlers:     map[int]DamageModHandler{},
		resolveMoveHandlers:     map[int]MoveHandler{},
		onHitAbilityHandlers:    map[int]OnHitHandler{},
		onHitItemHandlers:       map[int]OnHitHandler{},
		postMoveAbilityHandlers: map[int]PostMoveHandler{},
		postMoveItemHandlers:    map[int]PostMoveHandler{},
	}

	return r
}

func (r *Registry) RegisterEntryAbility(abilityId int, h EntryHandler) {
	r.entryAbilityHandlers[abilityId] = h
}

func (r *Registry) RegisterEntryItem(itemId int, h EntryHandler) {
	r.entryItemHandlers[itemId] = h
}

func (r *Registry) RegisterExitAbility(abilityId int, h ExitHandler) {
	r.exitAbilityHandlers[abilityId] = h
}

func (r *Registry) RegisterExitItem(itemId int, h ExitHandler) {
	r.exitItemHandlers[itemId] = h
}

func (r *Registry) AddResolveBase(h DamageModHandler) {
	r.resolveBaseHandlers = append(r.resolveBaseHandlers, h)
}

func (r *Registry) RegisterResolveAbility(abilityId int, h DamageModHandler) {
	r.resolveAbilityHandlers[abilityId] = h
}

func (r *Registry) RegisterResolveItem(itemId int, h DamageModHandler) {
	r.resolveItemHandlers[itemId] = h
}

func (r *Registry) RegisterResolveMove(moveId int, h MoveHandler) {
	r.resolveMoveHandlers[moveId] = h
}

func (r *Registry) RegisterOnHitAbility(abilityId int, h OnHitHandler) {
	r.onHitAbilityHandlers[abilityId] = h
}

func (r *Registry) RegisterOnHitItem(itemId int, h OnHitHandler) {
	r.onHitItemHandlers[itemId] = h
}

func (r *Registry) RegisterPostMoveAbility(abilityId int, h PostMoveHandler) {
	r.postMoveAbilityHandlers[abilityId] = h
}

func (r *Registry) RegisterPostMoveItem(itemId int, h PostMoveHandler) {
	r.postMoveItemHandlers[itemId] = h
}
