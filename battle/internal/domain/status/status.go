package status

type Status struct {
	main  *MainStatus
	other []OtherStatus
}

func NewStatus() Status {
	return Status{main: nil, other: nil}
}

func (s *Status) SetMainStatus(m *MainStatus) {
	if s.main != nil {
		return
	}
	s.main = m
}

func (s *Status) ForceSetMainStatus(m *MainStatus) {
	s.main = m
}

// TODO: Main/Other accessor は今後追加
