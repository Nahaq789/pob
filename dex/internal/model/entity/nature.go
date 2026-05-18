package entity

type Nature struct {
	Id              int    `gorm:"column:id;primaryKey"`
	Name            string `gorm:"column:name"`
	IncreasedStatId *int   `gorm:"column:increased_stat_id"` // NULL: ステータス補正なし（例: てれや）
	DecreasedStatId *int   `gorm:"column:decreased_stat_id"` // NULL: ステータス補正なし（例: てれや）
}

func (Nature) TableName() string { return "natures" }
