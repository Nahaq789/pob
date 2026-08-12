package pokemon

import (
	"pob/battle/internal/domain/ability"
	"pob/battle/internal/domain/apperror"
	"pob/battle/internal/domain/hp"
	"pob/battle/internal/domain/item"
	"pob/battle/internal/domain/move"
	"pob/battle/internal/domain/nature"
	"pob/battle/internal/domain/ptype"
	"pob/battle/internal/domain/rank"
	"pob/battle/internal/domain/status"
	"pob/battle/internal/domain/status/other"
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
	name      string
	slot      int
	nickname  string
	types     [2]ptype.Type
	baseStats BaseStats
	realStats RealStats
	nature    nature.Nature
	ability   *ability.Ability
	moves     [4]*move.Move

	// 動的データ
	currentHP          hp.HP
	ranks              rank.Rank
	mainStatus         *status.MainStatus
	otherStatuses      []other.OtherStatus
	heldItem           *item.Item
	lastConsumedItem   *item.Item
	lastSelectedMoveId int
	// このターンに場に出たばかりかフラグ
	justEntered bool
	events      []DomainEvent
}

// NewPokemon はPokemonのコンストラクタ。
// box-serviceから集約された全データを一括で受け取る想定のため、
// 部分的な省略を許容しないフルコンストラクタとする。
func NewPokemon(
	id PokemonId,
	name string,
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
	otherStatuses []other.OtherStatus,
	heldItem *item.Item,
	lastConsumedItem *item.Item,
	justEntered bool,
) *Pokemon {
	return &Pokemon{
		id:               id,
		name:             name,
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

func (p *Pokemon) Exited() {
	p.justEntered = false
}

func (p *Pokemon) IsFainted() bool {
	return p.currentHP.IsEmpty()
}

func (p *Pokemon) ResetOnSwitchOut() {
}

// 素早さの実数値
func (p *Pokemon) Speed() int {
	s := float64(p.realStats.Speed) * p.ranks.Speed().Value()
	// 状態異常がまひの場合は素早さが1/4になる
	if p.mainStatus != nil && p.mainStatus.Condition() == status.Paralysis {
		s *= 0.25
	}
	return int(s)
}

func (p *Pokemon) Id() PokemonId { return p.id }

func (p *Pokemon) Name() string { return p.name }

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

func (p *Pokemon) Moves() [4]*move.Move { return p.moves }
func (p *Pokemon) MoveById(moveId int) (*move.Move, error) {
	for _, m := range p.moves {
		if m == nil {
			continue
		}
		if m.Id() == moveId {
			return m, nil
		}
	}
	return nil, apperror.ErrMoveNotFound
}

// func (p *Pokemon) CurrentHP() vo.Count { return p.currentHP }
//
// func (p *Pokemon) Ranks() rank.Rank { return p.ranks }
//
// func (p *Pokemon) MainStatus() status.MainStatus { return p.mainStatus }
//
// func (p *Pokemon) OtherStatuses() []statusother.OtherStatus { return p.otherStatuses }
func (p *Pokemon) HeldItem() *item.Item { return p.heldItem }

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

func (p *Pokemon) LastSelectedMoveId() int { return p.lastSelectedMoveId }

func (p *Pokemon) SetLastSelectedMoveId(moveId int) {
	p.lastSelectedMoveId = moveId
}
