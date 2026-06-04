package entity

type Move struct {
	Id          int    `gorm:"column:id;primaryKey"`
	Name        string `gorm:"column:name"`
	TypeId      int    `gorm:"column:type_id"`
	DamageClass string `gorm:"column:damage_class"`
	Power       *int   `gorm:"column:power"`    // NULL: 固定・変動ダメージ
	Accuracy    *int   `gorm:"column:accuracy"` // NULL: 必中
	Pp          int    `gorm:"column:pp"`
	Priority    int    `gorm:"column:priority"`
}

func (Move) TableName() string { return "moves" }
