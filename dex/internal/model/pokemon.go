package model

type PokemonList struct {
	PokemonId int
	Name      string
}

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
	WeightKg      float64
	Abilities     []AbilitySlot
}

type AbilitySlot struct {
	AbilityId   int
	AbilityName string
	IsHidden    bool
}
