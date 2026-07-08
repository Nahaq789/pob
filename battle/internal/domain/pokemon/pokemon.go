package pokemon

import (
	"pob/battle/internal/domain/ability"
	"pob/battle/internal/domain/item"
	"pob/battle/internal/domain/move"
	"pob/battle/internal/domain/pp"
	"pob/battle/internal/domain/ptype"
	"pob/battle/internal/domain/rank"
	"pob/battle/internal/domain/status"
	"pob/battle/internal/domain/vo"
)

type PokemonId int

// 種族値を管理する専用パッケージ実装後に差し替え予定。
type BaseStats struct {
	HP, Attack, Defense, SpAttack, SpDefense, Speed int
}

// pkg/stats（実数値計算）と対応する専用パッケージ実装後に差し替え予定。
type RealStats struct {
	HP, Attack, Defense, SpAttack, SpDefense, Speed int
}

// 性格補正を扱う専用パッケージ実装後に差し替え予定。
type Nature string

// Pokemon はバトル中の1体のポケモンを表す集約。
// フィールドが多いため、値のコピーを避けるためポインタレシーバーで統一する。
type Pokemon struct {
	// 固定データ
	id        PokemonId
	slot      int
	nickname  string
	types     [2]ptype.Type
	baseStats BaseStats
	realStats RealStats
	nature    Nature
	ability   *ability.Ability
	moves     [4]*move.Move

	// 動的データ
	currentHP        vo.Count
	ranks            rank.Rank
	mainStatus       status.MainStatus
	otherStatuses    []status.OtherStatus
	pp               [4]pp.PP
	heldItem         *item.Item
	lastConsumedItem *item.Item
}

// NewPokemon はPokemonのコンストラクタ。
// box-serviceから集約された全データを一括で受け取る想定のため、
// 部分的な省略を許容しないフルコンストラクタとする。
func NewPokemon(
	id PokemonId,
	slot int,
	nickname string,
	types [2]ptype.Type,
	baseStats BaseStats,
	realStats RealStats,
	nature Nature,
	ability *ability.Ability,
	moves [4]*move.Move,
	currentHP vo.Count,
	ranks rank.Rank,
	mainStatus status.MainStatus,
	otherStatuses []status.OtherStatus,
	pp [4]pp.PP,
	heldItem *item.Item,
	lastConsumedItem *item.Item,
) *Pokemon {
	return &Pokemon{
		id:               id,
		slot:             slot,
		nickname:         nickname,
		types:            types,
		baseStats:        baseStats,
		realStats:        realStats,
		nature:           nature,
		ability:          ability,
		moves:            moves,
		currentHP:        currentHP,
		ranks:            ranks,
		mainStatus:       mainStatus,
		otherStatuses:    otherStatuses,
		pp:               pp,
		heldItem:         heldItem,
		lastConsumedItem: lastConsumedItem,
	}
}

func (p *Pokemon) ID() PokemonId { return p.id }

func (p *Pokemon) Slot() int { return p.slot }

func (p *Pokemon) Nickname() string { return p.nickname }

func (p *Pokemon) Types() [2]ptype.Type { return p.types }

func (p *Pokemon) BaseStats() BaseStats { return p.baseStats }

func (p *Pokemon) RealStats() RealStats { return p.realStats }

func (p *Pokemon) Nature() Nature { return p.nature }

func (p *Pokemon) Ability() *ability.Ability { return p.ability }

func (p *Pokemon) Moves() [4]*move.Move { return p.moves }

func (p *Pokemon) CurrentHP() vo.Count { return p.currentHP }

func (p *Pokemon) Ranks() rank.Rank { return p.ranks }

func (p *Pokemon) MainStatus() status.MainStatus { return p.mainStatus }

func (p *Pokemon) OtherStatuses() []status.OtherStatus { return p.otherStatuses }

func (p *Pokemon) PP() [4]pp.PP { return p.pp }

func (p *Pokemon) HeldItem() *item.Item { return p.heldItem }

func (p *Pokemon) LastConsumedItem() *item.Item { return p.lastConsumedItem }
