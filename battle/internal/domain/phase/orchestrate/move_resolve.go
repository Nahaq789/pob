package orchestrate

type MoveResolveOrchestrator struct {
}

func NewMoveResolveOrchestrator() *MoveResolveOrchestrator {
	return &MoveResolveOrchestrator{}
}

func (m *MoveResolveOrchestrator) Handle() (map[string][]string, error) {
	return nil, nil
}
