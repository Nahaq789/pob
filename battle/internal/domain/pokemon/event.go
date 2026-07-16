package pokemon

type EventKind string

type DomainEvent struct {
	Kind      EventKind
	PokemonId PokemonId
}

const (
	EventEntered EventKind = "entered"
)

func (p *Pokemon) PullEvent() []DomainEvent {
	events := p.events
	p.events = nil
	return events
}
