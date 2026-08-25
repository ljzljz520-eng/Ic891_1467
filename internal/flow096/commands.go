package flow096

import (
	"coursecodes/internal/domain"
	"coursecodes/internal/store"
)

type Command struct {
	RecordID string
	Quantity int
	Actor    string
}

func ApplyChange(s *store.Store, c Command) (domain.Record, error) {
	r, e := s.GetRecord(c.RecordID)
	if e != nil {
		return r, e
	}
	if e = r.ChangeQuantity(c.Quantity); e != nil {
		return r, e
	}
	return r, s.SaveRecord(r)
}
func ValidateCommand(c Command) bool         { return c.RecordID != "" && c.Quantity >= 0 && c.Actor != "" }
func CompareVersion(a, b domain.Record) bool { return a.ID == b.ID && a.Version <= b.Version }
