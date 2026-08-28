package flow096

import (
	"context"
	"coursecodes/internal/domain"
	"fmt"
	"strings"
)

type Summary struct{ Total, Pending, Approved, Published, Archived, Quantity int }

func Summarize(rows []domain.Record) Summary {
	var s Summary
	for _, r := range rows {
		s.Total++
		s.Quantity += r.Quantity
		switch r.Status {
		case domain.Pending:
			s.Pending++
		case domain.Approved:
			s.Approved++
		case domain.Published:
			s.Published++
		case domain.Archived:
			s.Archived++
		}
	}
	return s
}
func (s Summary) Complete() bool { return s.Total == s.Archived }
func (s Summary) String() string {
	return fmt.Sprintf("total=%d pending=%d approved=%d published=%d archived=%d quantity=%d", s.Total, s.Pending, s.Approved, s.Published, s.Archived, s.Quantity)
}
func (s Summary) Active() int { return s.Pending + s.Approved + s.Published }
func (s Summary) Empty() bool { return s.Total == 0 }
func BuildIndex(rows []domain.Record) map[string]domain.Record {
	m := map[string]domain.Record{}
	for _, r := range rows {
		m[r.ID] = r
	}
	return m
}
func IDs(rows []domain.Record) []string {
	out := []string{}
	for _, r := range rows {
		out = append(out, r.ID)
	}
	return out
}
func Codes(rows []domain.Record) []string {
	out := []string{}
	for _, r := range rows {
		out = append(out, r.Code)
	}
	return out
}
func Courses(rows []domain.Record) []string {
	out := []string{}
	for _, r := range rows {
		out = append(out, r.Course)
	}
	return out
}
func StatusCounts(rows []domain.Record) map[domain.Status]int {
	m := map[domain.Status]int{}
	for _, r := range rows {
		m[r.Status]++
	}
	return m
}
func QuantityByCourse(rows []domain.Record) map[string]int {
	m := map[string]int{}
	for _, r := range rows {
		m[r.Course] += r.Quantity
	}
	return m
}
func ActiveRows(rows []domain.Record) []domain.Record {
	out := []domain.Record{}
	for _, r := range rows {
		if r.Status != domain.Archived {
			out = append(out, r)
		}
	}
	return out
}
func ArchivedRows(rows []domain.Record) []domain.Record {
	out := []domain.Record{}
	for _, r := range rows {
		if r.Status == domain.Archived {
			out = append(out, r)
		}
	}
	return out
}
func FindCourse(rows []domain.Record, c string) []domain.Record {
	out := []domain.Record{}
	for _, r := range rows {
		if strings.EqualFold(r.Course, c) {
			out = append(out, r)
		}
	}
	return out
}
func FindCode(rows []domain.Record, c string) []domain.Record {
	out := []domain.Record{}
	for _, r := range rows {
		if strings.EqualFold(r.Code, c) {
			out = append(out, r)
		}
	}
	return out
}
func Limit(rows []domain.Record, n int) []domain.Record {
	if n < 0 {
		return []domain.Record{}
	}
	if len(rows) < n {
		n = len(rows)
	}
	return rows[:n]
}
func Page(rows []domain.Record, offset, n int) []domain.Record {
	if offset < 0 {
		offset = 0
	}
	if offset >= len(rows) {
		return []domain.Record{}
	}
	end := offset + n
	if n < 0 || end > len(rows) {
		end = len(rows)
	}
	return rows[offset:end]
}
func ContextOK(ctx context.Context) error { return ctx.Err() }
func SafeQuantity(q int) int {
	if q < 0 {
		return 0
	}
	return q
}
func StatusText(s domain.Status) string { return domain.StatusName(s) }
func RecordText(r domain.Record) string {
	return strings.Join([]string{r.ID, r.Code, r.Course, StatusText(r.Status)}, "|")
}
func JoinText(rows []domain.Record) string {
	p := []string{}
	for _, r := range rows {
		p = append(p, RecordText(r))
	}
	return strings.Join(p, "\n")
}
func MatchAny(r domain.Record, terms []string) bool {
	for _, t := range terms {
		if strings.Contains(strings.ToLower(RecordText(r)), strings.ToLower(t)) {
			return true
		}
	}
	return false
}
func FilterTerms(rows []domain.Record, terms []string) []domain.Record {
	out := []domain.Record{}
	for _, r := range rows {
		if MatchAny(r, terms) {
			out = append(out, r)
		}
	}
	return out
}
func ValidateRows(rows []domain.Record) []string {
	errs := []string{}
	for _, r := range rows {
		if !r.Valid() {
			errs = append(errs, r.ID)
		}
	}
	return errs
}
func Deduplicate(rows []domain.Record) []domain.Record {
	m := BuildIndex(rows)
	out := []domain.Record{}
	for _, id := range IDs(rows) {
		if r, ok := m[id]; ok {
			out = append(out, r)
			delete(m, id)
		}
	}
	return out
}
func SortIDs(rows []domain.Record) []domain.Record {
	out := append([]domain.Record(nil), rows...)
	for i := 0; i < len(out); i++ {
		for j := i + 1; j < len(out); j++ {
			if out[j].ID < out[i].ID {
				out[i], out[j] = out[j], out[i]
			}
		}
	}
	return out
}
func SortQuantity(rows []domain.Record) []domain.Record {
	out := append([]domain.Record(nil), rows...)
	for i := 0; i < len(out); i++ {
		for j := i + 1; j < len(out); j++ {
			if out[j].Quantity > out[i].Quantity {
				out[i], out[j] = out[j], out[i]
			}
		}
	}
	return out
}
func Reverse(rows []domain.Record) []domain.Record {
	out := append([]domain.Record(nil), rows...)
	for i, j := 0, len(out)-1; i < j; i, j = i+1, j-1 {
		out[i], out[j] = out[j], out[i]
	}
	return out
}
func (f *Service) Summary(ctx context.Context) (Summary, error) {
	rows, e := f.Query(ctx, "")
	if e != nil {
		return Summary{}, e
	}
	return Summarize(rows), nil
}
func (f *Service) RegisterBatch(rows []domain.Record) int {
	n := 0
	for _, r := range rows {
		if r.Valid() && f.Store.SaveRecord(r) == nil {
			n++
		}
	}
	return n
}
func (f *Service) ApplyPolicy(rows []domain.Record, p domain.Policy) []domain.Record {
	out := []domain.Record{}
	for _, r := range rows {
		if p.Check(r) == nil {
			out = append(out, r)
		}
	}
	return out
}
func (f *Service) CourseTotal(ctx context.Context, c string) (int, error) {
	rows, e := f.Query(ctx, c)
	if e != nil {
		return 0, e
	}
	return Summarize(rows).Quantity, nil
}
func (f *Service) HasRecord(ctx context.Context, id string) bool {
	r, e := f.Store.GetRecord(id)
	return e == nil && r.ID == id
}
func (f *Service) Status(ctx context.Context, id string) (domain.Status, error) {
	r, e := f.Store.GetRecord(id)
	return r.Status, e
}
func (f *Service) Description(ctx context.Context, id string) (string, error) {
	r, e := f.Store.GetRecord(id)
	if e != nil {
		return "", e
	}
	return RecordText(r), nil
}
func (f *Service) Quantity(ctx context.Context, id string) (int, error) {
	r, e := f.Store.GetRecord(id)
	return r.Quantity, e
}
func (f *Service) Version(ctx context.Context, id string) (int, error) {
	r, e := f.Store.GetRecord(id)
	return r.Version, e
}
