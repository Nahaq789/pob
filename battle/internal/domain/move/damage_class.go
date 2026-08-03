package move

type DamageClass string

const (
	DamageClassPhysical DamageClass = "physical"
	DamageClassSpecial  DamageClass = "special"
	DamageClassStatus   DamageClass = "status"
)
