package store

import (
	"coursecodes/internal/domain"
	"encoding/json"
	"errors"
	"go.etcd.io/bbolt"
	"path/filepath"
	"sort"
)

var recordsBucket = []byte("records")
var eventsBucket = []byte("events")
var workflowsBucket = []byte("workflows")
var attachmentsBucket = []byte("attachments")

type Store struct{ db *bbolt.DB }

func Open(path string) (*Store, error) {
	db, err := bbolt.Open(filepath.Clean(path), 0600, nil)
	if err != nil {
		return nil, err
	}
	s := &Store{db: db}
	err = db.Update(func(tx *bbolt.Tx) error {
		for _, b := range [][]byte{recordsBucket, eventsBucket, workflowsBucket, attachmentsBucket} {
			if _, e := tx.CreateBucketIfNotExists(b); e != nil {
				return e
			}
		}
		return nil
	})
	if err != nil {
		db.Close()
		return nil, err
	}
	return s, nil
}
func (s *Store) Close() error {
	if s == nil || s.db == nil {
		return nil
	}
	return s.db.Close()
}
func put(tx *bbolt.Tx, b []byte, key string, v any) error {
	x, e := json.Marshal(v)
	if e != nil {
		return e
	}
	return tx.Bucket(b).Put([]byte(key), x)
}
func get(tx *bbolt.Tx, b []byte, key string, v any) error {
	x := tx.Bucket(b).Get([]byte(key))
	if x == nil {
		return errors.New("not found")
	}
	return json.Unmarshal(x, v)
}
func (s *Store) SaveRecord(r domain.Record) error {
	return s.db.Update(func(tx *bbolt.Tx) error { return put(tx, recordsBucket, r.ID, r) })
}
func (s *Store) GetRecord(id string) (domain.Record, error) {
	var r domain.Record
	e := s.db.View(func(tx *bbolt.Tx) error { return get(tx, recordsBucket, id, &r) })
	return r, e
}
func (s *Store) ListRecords() ([]domain.Record, error) {
	out := []domain.Record{}
	e := s.db.View(func(tx *bbolt.Tx) error {
		return tx.Bucket(recordsBucket).ForEach(func(_, v []byte) error {
			var r domain.Record
			if e := json.Unmarshal(v, &r); e != nil {
				return e
			}
			out = append(out, r)
			return nil
		})
	})
	sort.Slice(out, func(i, j int) bool { return out[i].ID < out[j].ID })
	return out, e
}
func (s *Store) SaveEvent(e domain.AuditEvent) error {
	return s.db.Update(func(tx *bbolt.Tx) error { return put(tx, eventsBucket, e.ID, e) })
}
func (s *Store) ListEvents(recordID string) ([]domain.AuditEvent, error) {
	out := []domain.AuditEvent{}
	e := s.db.View(func(tx *bbolt.Tx) error {
		return tx.Bucket(eventsBucket).ForEach(func(_, v []byte) error {
			var a domain.AuditEvent
			if x := json.Unmarshal(v, &a); x != nil {
				return x
			}
			if recordID == "" || a.RecordID == recordID {
				out = append(out, a)
			}
			return nil
		})
	})
	return out, e
}
func (s *Store) SaveWorkflow(w domain.Workflow) error {
	return s.db.Update(func(tx *bbolt.Tx) error { return put(tx, workflowsBucket, w.ID, w) })
}
func (s *Store) SaveAttachment(a domain.Attachment) error {
	return s.db.Update(func(tx *bbolt.Tx) error { return put(tx, attachmentsBucket, a.ID, a) })
}
func (s *Store) DBPath() string { return s.db.Path() }
