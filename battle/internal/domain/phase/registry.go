package phase

type Registry struct {
	entryAbilityHandlers map[int]EntryHandler
	entryItemHandlers    map[int]EntryHandler

	exitAbilityHandlers map[int]ExitHandler
	exitItemHandlers    map[int]ExitHandler

	damageBaseHandlers    []DamageHandler
	damageAbilityHandlers map[int]DamageHandler
	damageItemHandlers    map[int]DamageHandler
	damageMoveHandlers    map[int]DamageHandler
}

func NewRegistry() *Registry {
	r := &Registry{
		entryAbilityHandlers:  map[int]EntryHandler{},
		entryItemHandlers:     map[int]EntryHandler{},
		exitAbilityHandlers:   map[int]ExitHandler{},
		exitItemHandlers:      map[int]ExitHandler{},
		damageBaseHandlers:    []DamageHandler{},
		damageAbilityHandlers: map[int]DamageHandler{},
		damageItemHandlers:    map[int]DamageHandler{},
		damageMoveHandlers:    map[int]DamageHandler{},
	}

	return r
}
