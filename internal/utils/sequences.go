package utils

type sequencer struct {
	value int
}

func (s sequencer) NextVal() int {
	s.value++
	return s.value
}

func NewSequence() *sequencer {
	return new(sequencer)
}
