package dto

type CreatePartyRequest struct {
	Name string `json:"name"`
}

type UpdatePartyNameRequest struct {
	Name string `json:"name"`
}

type PartyPokemonEntry struct {
	BoxPokemonId string `json:"box_pokemon_id"`
	Slot         int    `json:"slot"`
}

type SetPartyPokemonRequest struct {
	Pokemon []PartyPokemonEntry `json:"pokemon"`
}
