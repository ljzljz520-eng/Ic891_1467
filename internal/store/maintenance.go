package store

import (
	"coursecodes/internal/domain"
	"fmt"
	"go.etcd.io/bbolt"
)

func (s *Store) EnsureVersion(id string, want int) error {
	r, e := s.GetRecord(id)
	if e != nil {
		return e
	}
	if r.Version < want {
		r.Version = want
		return s.SaveRecord(r)
	}
	return nil
}
func (s *Store) DeleteArchived() int {
	rows, e := s.ListRecords()
	if e != nil {
		return 0
	}
	n := 0
	for _, r := range rows {
		if r.Status == domain.Archived {
			_ = s.db.Update(func(tx *bbolt.Tx) error { return tx.Bucket(recordsBucket).Delete([]byte(r.ID)) })
			n++
		}
	}
	return n
}
func Key(prefix, id string) string { return fmt.Sprintf("%s:%s", prefix, id) }
