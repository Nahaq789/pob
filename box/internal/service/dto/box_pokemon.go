package dto

type AddBoxPokemonRequest struct {
	PokemonId  int
	Nickname   *string
	AbilityId  int
	Nature     string
	HeldItemId *int
	IvHp       int
	IvAttack   int
	IvDefense  int
	IvSpAttack int
	IvSpDefense int
	IvSpeed    int
	EvHp       int
	EvAttack   int
	EvDefense  int
	EvSpAttack int
	EvSpDefense int
	EvSpeed    int
	Move1Id    *int
	Move2Id    *int
	Move3Id    *int
	Move4Id    *int
}

type UpdateBoxPokemonRequest struct {
	BoxPokemonId string
	Nickname     *string
	AbilityId    int
	Nature       string
	HeldItemId   *int
	IvHp         int
	IvAttack     int
	IvDefense    int
	IvSpAttack   int
	IvSpDefense  int
	IvSpeed      int
	EvHp         int
	EvAttack     int
	EvDefense    int
	EvSpAttack   int
	EvSpDefense  int
	EvSpeed      int
	Move1Id      *int
	Move2Id      *int
	Move3Id      *int
	Move4Id      *int
}
