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

	postMoveAbilityHandlers map[int]PostMoveHandler
	postMoveItemHandlers    map[int]PostMoveHandler
}

func NewRegistry() *Registry {
	r := &Registry{
		entryAbilityHandlers:   map[int]EntryHandler{},
		entryItemHandlers:      map[int]EntryHandler{},
		exitAbilityHandlers:    map[int]ExitHandler{},
		exitItemHandlers:       map[int]ExitHandler{},
		resolveBaseHandlers:    []DamageModHandler{},
		resolveAbilityHandlers: map[int]DamageModHandler{},
		resolveItemHandlers:    map[int]DamageModHandler{},
		resolveMoveHandlers:     map[int]MoveHandler{},
		postMoveAbilityHandlers: map[int]PostMoveHandler{},
		postMoveItemHandlers:    map[int]PostMoveHandler{},
	}

	return r
}
