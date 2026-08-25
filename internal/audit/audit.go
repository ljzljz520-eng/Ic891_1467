package audit

import (
	"coursecodes/internal/domain"
	"coursecodes/internal/store"
	"fmt"
)

type Logger struct {
	Store *store.Store
	seq   int
}

func New(s *store.Store) *Logger { return &Logger{Store: s} }
func (l *Logger) Record(id, action, actor string) error {
	l.seq++
	e := domain.AuditEvent{ID: fmt.Sprintf("event-%06d", l.seq), RecordID: id, Action: action, Actor: actor, At: "deterministic"}
	return l.Store.SaveEvent(e)
}
func (l *Logger) ForRecord(id string) ([]domain.AuditEvent, error) { return l.Store.ListEvents(id) }
func Summarize(es []domain.AuditEvent) map[string]int {
	m := map[string]int{}
	for _, e := range es {
		m[e.Action]++
	}
	return m
}
func HasAction(es []domain.AuditEvent, a string) bool {
	for _, e := range es {
		if e.Action == a {
			return true
		}
	}
	return false
}
