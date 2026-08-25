package store

import (
	"coursecodes/internal/domain"
	"encoding/json"
)

func (s *Store) Snapshot() ([]byte, error) {
	rows, e := s.ListRecords()
	if e != nil {
		return nil, e
	}
	return json.Marshal(rows)
}
func (s *Store) Restore(data []byte) error {
	var rows []domain.Record
	if e := json.Unmarshal(data, &rows); e != nil {
		return e
	}
	for _, r := range rows {
		if e := s.SaveRecord(r); e != nil {
			return e
		}
	}
	return nil
}
func (s *Store) Count() int {
	r, e := s.ListRecords()
	if e != nil {
		return 0
	}
	return len(r)
}
