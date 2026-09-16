package phase

type Registry struct {
	entryAbilityHandlers map[int]EntryHandler
	entryItemHandlers    map[int]EntryHandler

	exitAbilityHandlers map[int]ExitHandler
	exitItemHandlers    map[int]ExitHandler

	resolveBaseHandlers    []DamageModHandler
	resolveAbilityHandlers map[int]DamageModHandler
	resolveItemHandlers    map[int]DamageModHandler
	resolveMoveHandlers map[int]MoveHandler
}

func NewRegistry() *Registry {
	r := &Registry{
		entryAbilityHandlers:  map[int]EntryHandler{},
		entryItemHandlers:     map[int]EntryHandler{},
		exitAbilityHandlers:   map[int]ExitHandler{},
		exitItemHandlers:      map[int]ExitHandler{},
		resolveBaseHandlers:    []DamageModHandler{},
		resolveAbilityHandlers: map[int]DamageModHandler{},
		resolveItemHandlers:    map[int]DamageModHandler{},
		resolveMoveHandlers: map[int]MoveHandler{},
	}

	return r
}
