package phase

type Registry struct {
	entryAbilityHandler map[int]EntryHandler
	entryItemHandler    map[int]EntryHandler
}

func NewRegistry() *Registry {
	r := &Registry{
		entryAbilityHandler: map[int]EntryHandler{},
		entryItemHandler:    map[int]EntryHandler{},
	}

	return r
}
