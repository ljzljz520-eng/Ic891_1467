package flow096

import (
	"context"
	"coursecodes/internal/audit"
	"coursecodes/internal/domain"
	"coursecodes/internal/search"
	"coursecodes/internal/store"
	"errors"
)

type Service struct {
	Store  *store.Store
	Search *search.Service
	Audit  *audit.Logger
}

func New(s *store.Store) *Service {
	return &Service{Store: s, Search: search.New(s), Audit: audit.New(s)}
}
func (f *Service) Register(id, code, course string, q int) (domain.Record, error) {
	r, e := domain.NewRecord(id, code, course, q)
	if e != nil {
		return r, e
	}
	e = f.Store.SaveRecord(r)
	if e == nil {
		e = f.Audit.Record(id, "register", "operator")
	}
	return r, e
}
func (f *Service) Review(id string, approve bool) error {
	r, e := f.Store.GetRecord(id)
	if e != nil {
		return e
	}
	if !approve {
		return errors.New("review rejected")
	}
	if e = r.Approve(); e != nil {
		return e
	}
	if e = f.Store.SaveRecord(r); e != nil {
		return e
	}
	return f.Audit.Record(id, "approve", "reviewer")
}
func (f *Service) Change(id string, q int) error {
	r, e := f.Store.GetRecord(id)
	if e != nil {
		return e
	}
	if e = r.ChangeQuantity(q); e != nil {
		return e
	}
	if e = f.Store.SaveRecord(r); e != nil {
		return e
	}
	return f.Audit.Record(id, "change", "operator")
}
func (f *Service) Publish(id string) error {
	r, e := f.Store.GetRecord(id)
	if e != nil {
		return e
	}
	if e = r.Publish(); e != nil {
		return e
	}
	if e = f.Store.SaveRecord(r); e != nil {
		return e
	}
	return f.Audit.Record(id, "publish", "operator")
}
func (f *Service) Archive(id string) error {
	r, e := f.Store.GetRecord(id)
	if e != nil {
		return e
	}
	if e = r.Archive(); e != nil {
		return e
	}
	if e = f.Store.SaveRecord(r); e != nil {
		return e
	}
	return f.Audit.Record(id, "archive", "operator")
}
func (f *Service) Query(ctx context.Context, q string) ([]domain.Record, error) {
	return f.Search.Find(ctx, q)
}
func (f *Service) Workflow(ctx context.Context, id string) error {
	if err := ctx.Err(); err != nil {
		return err
	}
	if e := f.Review(id, true); e != nil {
		return e
	}
	if e := f.Publish(id); e != nil {
		return e
	}
	return f.Archive(id)
}
