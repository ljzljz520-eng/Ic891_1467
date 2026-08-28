package domain

import (
	"errors"
	"strings"
)

type Status string

const (
	Pending   Status = "pending"
	Approved  Status = "approved"
	Published Status = "published"
	Archived  Status = "archived"
)

type Record struct {
	ID, Code, Course string
	Quantity         int
	Status           Status
	Version          int
}
type AuditEvent struct{ ID, RecordID, Action, Actor, At string }
type Workflow struct{ ID, RecordID, Stage, State string }
type Attachment struct{ ID, RecordID, Name, Digest string }

func NewRecord(id, code, course string, quantity int) (Record, error) {
	if strings.TrimSpace(id) == "" || strings.TrimSpace(code) == "" {
		return Record{}, errors.New("id and code required")
	}
	if quantity < 0 {
		return Record{}, errors.New("quantity must be nonnegative")
	}
	return Record{ID: id, Code: code, Course: course, Quantity: quantity, Status: Pending, Version: 1}, nil
}
func (r Record) Review(approve bool) error {
	if r.Status != Pending {
		return errors.New("record not pending")
	}
	if !approve {
		return errors.New("review rejected")
	}
	return nil
}
func (r *Record) Approve() error {
	if err := r.Review(true); err != nil {
		return err
	}
	r.Status = Approved
	r.Version++
	return nil
}
func (r *Record) Publish() error {
	if r.Status != Approved {
		return errors.New("record not approved")
	}
	r.Status = Published
	r.Version++
	return nil
}
func (r *Record) Archive() error {
	if r.Status != Published {
		return errors.New("record not published")
	}
	r.Status = Archived
	r.Version++
	return nil
}
func (r *Record) ChangeQuantity(q int) error {
	if q < 0 {
		return errors.New("quantity must be nonnegative")
	}
	if r.Status == Archived {
		return errors.New("archived record")
	}
	r.Quantity = q
	r.Version++
	return nil
}
func (r Record) Valid() bool        { return r.ID != "" && r.Code != "" && r.Quantity >= 0 }
func NormalizeCode(s string) string { return strings.ToUpper(strings.TrimSpace(s)) }
func StatusName(s Status) string {
	switch s {
	case Pending:
		return "pending"
	case Approved:
		return "approved"
	case Published:
		return "published"
	case Archived:
		return "archived"
	default:
		return "unknown"
	}
}
