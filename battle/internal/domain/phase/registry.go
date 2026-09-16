package phase

type Registry struct {
	entryAbilityHandlers map[int]EntryHandler
	entryItemHandlers    map[int]EntryHandler

	exitAbilityHandlers map[int]ExitHandler
	exitItemHandlers    map[int]ExitHandler

	resolveBaseHandlers    []MoveResolveHandler
	resolveAbilityHandlers map[int]MoveResolveHandler
	resolveItemHandlers    map[int]MoveResolveHandler
	resolveMoveHandlers    map[int]MoveResolveHandler
}

func NewRegistry() *Registry {
	r := &Registry{
		entryAbilityHandlers:  map[int]EntryHandler{},
		entryItemHandlers:     map[int]EntryHandler{},
		exitAbilityHandlers:   map[int]ExitHandler{},
		exitItemHandlers:      map[int]ExitHandler{},
		resolveBaseHandlers:    []MoveResolveHandler{},
		resolveAbilityHandlers: map[int]MoveResolveHandler{},
		resolveItemHandlers:    map[int]MoveResolveHandler{},
		resolveMoveHandlers:    map[int]MoveResolveHandler{},
	}

	return r
}
