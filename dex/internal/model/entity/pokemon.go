package entity

type Pokemon struct {
	Id            int    `gorm:"column:id;primaryKey"`
	Name          string `gorm:"column:name"`
	Type1Id       int    `gorm:"column:type1_id"`
	Type2Id       *int   `gorm:"column:type2_id"` // NULL: 単タイプ
	BaseHp        int    `gorm:"column:base_hp"`
	BaseAttack    int    `gorm:"column:base_attack"`
	BaseDefense   int    `gorm:"column:base_defense"`
	BaseSpAttack  int    `gorm:"column:base_sp_attack"`
	BaseSpDefense int    `gorm:"column:base_sp_defense"`
	BaseSpeed     int    `gorm:"column:base_speed"`
}

func (Pokemon) TableName() string { return "pokemon" }
