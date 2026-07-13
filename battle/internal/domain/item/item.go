package item

type ItemId int
type ItemCategory string

type Item struct {
	id       ItemId
	category ItemCategory
}

const (
	CategoryBadHeldItems    ItemCategory = "bad-held-items"
	CategoryChoice          ItemCategory = "choice"
	CategoryMemories        ItemCategory = "memories"
	CategoryHeldItems       ItemCategory = "held-items"
	CategorySpeciesSpecific ItemCategory = "species-specific"
	CategoryMegaStones      ItemCategory = "mega-stones"
	CategoryMedicine        ItemCategory = "medicine"
	CategoryPlates          ItemCategory = "plates"
	CategoryTypeProtection  ItemCategory = "type-protection"
	CategoryJewels          ItemCategory = "jewels"
	CategoryTypeEnhancement ItemCategory = "type-enhancement"
	CategoryZCrystals       ItemCategory = "z-crystals"
	CategoryInAPinch        ItemCategory = "in-a-pinch"
	CategoryOther           ItemCategory = "other"
	CategoryPickyHealing    ItemCategory = "picky-healing"
)

var berryIds = map[ItemId]struct{}{
	126: {}, 127: {}, 128: {}, 129: {}, 130: {},
	131: {}, 132: {}, 133: {}, 134: {}, 135: {},
	136: {}, 137: {}, 138: {}, 139: {}, 140: {},
	161: {}, 162: {}, 163: {}, 164: {}, 165: {},
	166: {}, 167: {}, 168: {}, 169: {}, 170: {},
	171: {}, 172: {}, 173: {}, 174: {}, 175: {},
	176: {}, 177: {}, 178: {}, 179: {}, 180: {},
	181: {}, 182: {}, 183: {}, 184: {}, 185: {},
	186: {}, 187: {}, 188: {}, 189: {},
	723: {},
	724: {},
	725: {},
}

func NewItem(id ItemId, category ItemCategory) Item {
	return Item{id: id, category: category}
}

func (i Item) Id() ItemId             { return i.id }
func (i Item) Category() ItemCategory { return i.category }

func (i Item) IsBerry() bool {
	_, ok := berryIds[i.id]
	return ok
}

func (i Item) ConfuseFlavor() (Flavor, bool) {
	f, ok := confuseFlavors[i.id]
	return f, ok
}
