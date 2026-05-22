package model

type Pokemon struct {
	PokemonId     int
	Name          string
	Type1Id       int
	Type2Id       int
	BaseHp        int
	BaseAttack    int
	BaseDefense   int
	BaseSpAttack  int
	BaseSpDefense int
	BaseSpeed     int
	Abilities     []AbilitySlot
}

type AbilitySlot struct {
	AbilityId   int
	AbilityName string
	IsHidden    bool
}
