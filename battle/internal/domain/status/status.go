package status

type Status struct {
	main   *MainStatus
	others []OtherStatus
}

func NewStatus() Status {
	return Status{main: nil, others: []OtherStatus{}}
}

func NewStatusWith(main *MainStatus, others []OtherStatus) Status {
	if others == nil {
		others = []OtherStatus{}
	}
	return Status{main: main, others: others}
}

func (s Status) Main() *MainStatus { return s.main }

func (s Status) Others() []OtherStatus { return s.others }

func (s Status) OtherMap() map[OtherCondition]OtherStatus {
	m := make(map[OtherCondition]OtherStatus, len(s.others))
	for _, o := range s.others {
		m[o.Kind()] = o
	}
	return m
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

func (s *Status) AddOtherStatus(o OtherStatus) {
	s.others = append(s.others, o)
}

func (s *Status) RemoveOtherStatus(kind OtherCondition) {
	filtered := s.others[:0]
	for _, o := range s.others {
		if o.Kind() != kind {
			filtered = append(filtered, o)
		}
	}
	s.others = filtered
}
