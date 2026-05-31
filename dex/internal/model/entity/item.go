package entity

type Item struct {
	Id         int    `gorm:"column:id;primaryKey"`
	Name       string `gorm:"column:name"`
	Category   string `gorm:"column:category"`
	FlavorText string `gorm:"column:flavor_text"`
}

func (Item) TableName() string { return "items" }
