package entity

type PokemonMove struct {
	PokemonId int  `gorm:"column:pokemon_id;primaryKey"`
	MoveId    int  `gorm:"column:move_id;primaryKey"`
	Move      Move `gorm:"foreignKey:MoveId"`
}

func (PokemonMove) TableName() string { return "pokemon_moves" }
