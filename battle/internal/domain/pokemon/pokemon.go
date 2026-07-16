package pokemon

import (
	"pob/battle/internal/domain/ability"
	"pob/battle/internal/domain/hp"
	"pob/battle/internal/domain/item"
	"pob/battle/internal/domain/move"
	"pob/battle/internal/domain/nature"
	"pob/battle/internal/domain/pp"
	"pob/battle/internal/domain/ptype"
	"pob/battle/internal/domain/rank"
	"pob/battle/internal/domain/status"
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
	nature    nature.Nature
	ability   *ability.Ability
	moves     [4]*move.Move

	// 動的データ
	currentHP        hp.HP
	ranks            rank.Rank
	mainStatus       *status.MainStatus
	otherStatuses    []status.OtherStatus
	pp               [4]pp.PP
	heldItem         *item.Item
	lastConsumedItem *item.Item
	// このターンに場に出たばかりかフラグ
	justEntered bool
	events      []DomainEvent
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
	nature nature.Nature,
	ability *ability.Ability,
	moves [4]*move.Move,
	currentHP hp.HP,
	ranks rank.Rank,
	mainStatus *status.MainStatus,
	otherStatuses []status.OtherStatus,
	pp [4]pp.PP,
	heldItem *item.Item,
	lastConsumedItem *item.Item,
	justEntered bool,
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
		justEntered:      justEntered,
		events:           []DomainEvent{},
	}
}

// Entered はこのポケモンが場に出た直後の状態遷移をまとめる。
// 現状はjustEnteredのセットのみ。将来的に場に出た時の他の状態変化があればここに追加。
func (p *Pokemon) Entered() {
	p.justEntered = true
}

func (p *Pokemon) IsFainted() bool {
	return p.currentHP.IsEmpty()
}

func (p *Pokemon) ResetOnSwitchOut() {
}

func (p *Pokemon) Speed() int {
	return p.realStats.Speed
}

// func (p *Pokemon) ID() PokemonId { return p.id }
//
// func (p *Pokemon) Slot() int { return p.slot }
//
// func (p *Pokemon) Nickname() string { return p.nickname }
//
// func (p *Pokemon) Types() [2]ptype.Type { return p.types }
//
// func (p *Pokemon) BaseStats() BaseStats { return p.baseStats }
//
// func (p *Pokemon) RealStats() RealStats { return p.realStats }
//
// func (p *Pokemon) Nature() Nature { return p.nature }
func (p *Pokemon) Ability() *ability.Ability { return p.ability }

//
// func (p *Pokemon) Moves() [4]*move.Move { return p.moves }
//
// func (p *Pokemon) CurrentHP() vo.Count { return p.currentHP }
//
// func (p *Pokemon) Ranks() rank.Rank { return p.ranks }
//
// func (p *Pokemon) MainStatus() status.MainStatus { return p.mainStatus }
//
// func (p *Pokemon) OtherStatuses() []status.OtherStatus { return p.otherStatuses }
//
// func (p *Pokemon) PP() [4]pp.PP { return p.pp }
//
// func (p *Pokemon) HeldItem() *item.Item { return p.heldItem }
//
// func (p *Pokemon) LastConsumedItem() *item.Item { return p.lastConsumedItem }

// SetJustEntered は場に出た直後（初手選出・交代成立時）に呼び出す。
func (p *Pokemon) SetJustEntered() {
	p.justEntered = true
}

// ClearJustEntered はターン終了時（フェーズ6）に呼び出し、フラグを解除する。
func (p *Pokemon) ClearJustEntered() {
	p.justEntered = false
}

// IsJustEntered は場に出たばかりかどうかを返す。
func (p *Pokemon) IsJustEntered() bool {
	return p.justEntered
}
