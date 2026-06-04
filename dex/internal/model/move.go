package model

type Move struct {
	MoveId      int
	Name        string
	TypeId      int
	DamageClass string
	Power       int // 0: 固定・変動ダメージ
	Accuracy    int // 0: 必中
	Pp          int
	Priority    int
}
