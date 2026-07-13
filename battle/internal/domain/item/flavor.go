package item

type Flavor string

const (
	FlavorSpicy  Flavor = "spicy"  // からい
	FlavorSour   Flavor = "sour"   // すっぱい
	FlavorSweet  Flavor = "sweet"  // あまい
	FlavorBitter Flavor = "bitter" // にがい
	FlavorDry    Flavor = "dry"    // しぶい
)

var confuseFlavors = map[ItemId]Flavor{
	136: FlavorSpicy,  // フィラのみ
	137: FlavorDry,    // ウイのみ
	138: FlavorSweet,  // マゴのみ
	139: FlavorBitter, // バンジのみ
	140: FlavorSour,   // イアのみ
}
