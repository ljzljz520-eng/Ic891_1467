package search

import (
	"context"
	"coursecodes/internal/domain"
	"coursecodes/internal/store"
	"strings"
)

type Service struct{ Store *store.Store }

func New(s *store.Store) *Service { return &Service{Store: s} }
func (s *Service) Find(ctx context.Context, q string) ([]domain.Record, error) {
	if err := ctx.Err(); err != nil {
		return nil, err
	}
	all, e := s.Store.ListRecords()
	if e != nil {
		return nil, e
	}
	q = strings.ToUpper(strings.TrimSpace(q))
	out := []domain.Record{}
	for _, r := range all {
		if q == "" || strings.Contains(strings.ToUpper(r.Code), q) || strings.Contains(strings.ToUpper(r.Course), q) {
			out = append(out, r)
		}
	}
	return out, nil
}
func Match(r domain.Record, q string) bool {
	return q == "" || strings.Contains(strings.ToUpper(r.Code), strings.ToUpper(q)) || strings.Contains(strings.ToUpper(r.Course), strings.ToUpper(q))
}
func ByStatus(rs []domain.Record, st domain.Status) []domain.Record {
	out := []domain.Record{}
	for _, r := range rs {
		if r.Status == st {
			out = append(out, r)
		}
	}
	return out
}
func CountQuantity(rs []domain.Record) int {
	n := 0
	for _, r := range rs {
		n += r.Quantity
	}
	return n
}
