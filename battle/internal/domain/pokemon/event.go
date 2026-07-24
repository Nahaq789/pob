package pokemon

type EventKind string

type DomainEvent struct {
	Kind      EventKind
	PokemonId PokemonId
}

const (
	EventEntered EventKind = "entered" // ポケモンを出したときのイベント
)

func (p *Pokemon) PullEvents() []DomainEvent {
	events := p.events
	p.events = nil
	return events
}
