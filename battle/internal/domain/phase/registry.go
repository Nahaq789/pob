package phase

type Registry struct {
	entryAbilityHandlers map[int]EntryHandler
	entryItemHandlers    map[int]EntryHandler

	exitAbilityHandlers map[int]ExitHandler
	exitItemHandlers    map[int]ExitHandler
}

func NewRegistry() *Registry {
	r := &Registry{
		entryAbilityHandlers: map[int]EntryHandler{},
		entryItemHandlers:    map[int]EntryHandler{},
		exitAbilityHandlers:  map[int]ExitHandler{},
		exitItemHandlers:     map[int]ExitHandler{},
	}

	return r
}
