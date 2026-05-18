package entity

type Type struct {
	Id   int    `gorm:"column:id;primaryKey"`
	Name string `gorm:"column:name"`
}

func (Type) TableName() string { return "types" }
