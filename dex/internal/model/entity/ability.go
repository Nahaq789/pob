package entity

type Ability struct {
	Id   int    `gorm:"column:id;primaryKey"`
	Name string `gorm:"column:name"`
}

func (Ability) TableName() string { return "abilities" }
