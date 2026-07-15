package phase

type Registry struct {
	entryAbilityHandler map[int]EntryPhaseHandler
	entryItemHandler    map[int]EntryPhaseHandler
}

func NewRegistry() *Registry {
	r := &Registry{
		entryAbilityHandler: map[int]EntryPhaseHandler{},
		entryItemHandler:    map[int]EntryPhaseHandler{},
	}

	return r
}
