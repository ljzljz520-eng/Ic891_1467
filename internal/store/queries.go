package store

import (
	"coursecodes/internal/domain"
	"errors"
	"strings"
)

func (s *Store) FindCourse(course string) ([]domain.Record, error) {
	rows, e := s.ListRecords()
	if e != nil {
		return nil, e
	}
	out := []domain.Record{}
	for _, r := range rows {
		if strings.EqualFold(r.Course, course) {
			out = append(out, r)
		}
	}
	return out, nil
}
func (s *Store) FindCode(code string) ([]domain.Record, error) {
	rows, e := s.ListRecords()
	if e != nil {
		return nil, e
	}
	out := []domain.Record{}
	for _, r := range rows {
		if strings.EqualFold(r.Code, code) {
			out = append(out, r)
		}
	}
	return out, nil
}
func (s *Store) UpdateQuantity(id string, q int) error {
	r, e := s.GetRecord(id)
	if e != nil {
		return e
	}
	if e = r.ChangeQuantity(q); e != nil {
		return e
	}
	return s.SaveRecord(r)
}
func (s *Store) UpdateStatus(id string, status domain.Status) error {
	r, e := s.GetRecord(id)
	if e != nil {
		return e
	}
	r.Status = status
	r.Version++
	return s.SaveRecord(r)
}
func (s *Store) Exists(id string) bool { _, e := s.GetRecord(id); return e == nil }
func (s *Store) Quantity(id string) int {
	r, e := s.GetRecord(id)
	if e != nil {
		return 0
	}
	return r.Quantity
}
func (s *Store) Version(id string) int {
	r, e := s.GetRecord(id)
	if e != nil {
		return 0
	}
	return r.Version
}
func (s *Store) Status(id string) domain.Status {
	r, e := s.GetRecord(id)
	if e != nil {
		return ""
	}
	return r.Status
}
func (s *Store) ActiveCount() int {
	rows, e := s.ListRecords()
	if e != nil {
		return 0
	}
	n := 0
	for _, r := range rows {
		if r.Status != domain.Archived {
			n++
		}
	}
	return n
}
func (s *Store) ArchivedCount() int {
	rows, e := s.ListRecords()
	if e != nil {
		return 0
	}
	n := 0
	for _, r := range rows {
		if r.Status == domain.Archived {
			n++
		}
	}
	return n
}
func (s *Store) QuantityTotal() int {
	rows, e := s.ListRecords()
	if e != nil {
		return 0
	}
	n := 0
	for _, r := range rows {
		n += r.Quantity
	}
	return n
}
func (s *Store) Replace(r domain.Record) error {
	if !r.Valid() {
		return errors.New("invalid record")
	}
	return s.SaveRecord(r)
}
func (s *Store) RecordsByStatus(st domain.Status) ([]domain.Record, error) {
	rows, e := s.ListRecords()
	if e != nil {
		return nil, e
	}
	out := []domain.Record{}
	for _, r := range rows {
		if r.Status == st {
			out = append(out, r)
		}
	}
	return out, nil
}
func (s *Store) RecordsByQuantity(min int) ([]domain.Record, error) {
	rows, e := s.ListRecords()
	if e != nil {
		return nil, e
	}
	out := []domain.Record{}
	for _, r := range rows {
		if r.Quantity >= min {
			out = append(out, r)
		}
	}
	return out, nil
}
func (s *Store) IDs() []string {
	rows, e := s.ListRecords()
	if e != nil {
		return nil
	}
	out := []string{}
	for _, r := range rows {
		out = append(out, r.ID)
	}
	return out
}
func (s *Store) Codes() []string {
	rows, e := s.ListRecords()
	if e != nil {
		return nil
	}
	out := []string{}
	for _, r := range rows {
		out = append(out, r.Code)
	}
	return out
}
func (s *Store) Courses() []string {
	rows, e := s.ListRecords()
	if e != nil {
		return nil
	}
	out := []string{}
	for _, r := range rows {
		out = append(out, r.Course)
	}
	return out
}
