package model

type Gender int

const (
	GenderUnknown Gender = iota // 0
	GenderMale                  // 1
	GenderFemale                // 2
)

func (g Gender) Valid() bool {
	return g >= GenderUnknown && g <= GenderFemale
}
