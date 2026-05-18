package entity

type PokemonAbility struct {
	PokemonId int     `gorm:"column:pokemon_id;primaryKey"`
	AbilityId int     `gorm:"column:ability_id;primaryKey"`
	Slot      int     `gorm:"column:slot"`
	IsHidden  bool    `gorm:"column:is_hidden"`
	Ability   Ability `gorm:"foreignKey:AbilityId"`
}

func (PokemonAbility) TableName() string { return "pokemon_abilities" }
