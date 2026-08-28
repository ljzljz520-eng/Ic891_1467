package domain

import (
	"errors"
	"sort"
	"strings"
)

func (r Record) FieldMap() map[string]string {
	return map[string]string{"id": r.ID, "code": r.Code, "course": r.Course, "status": string(r.Status)}
}
func ValidateID(v string) bool     { return strings.TrimSpace(v) != "" && len(v) <= 64 }
func ValidateCode(v string) bool   { return len(strings.TrimSpace(v)) >= 3 && len(v) <= 64 }
func ValidateCourse(v string) bool { return strings.TrimSpace(v) != "" && len(v) <= 120 }
func QuantityBand(q int) string {
	if q < 0 {
		return "invalid"
	}
	if q == 0 {
		return "empty"
	}
	if q < 10 {
		return "small"
	}
	if q < 100 {
		return "medium"
	}
	return "large"
}
func IsMutable(s Status) bool { return s == Pending || s == Approved || s == Published }
func CloneRecord(r Record) Record {
	return Record{ID: r.ID, Code: r.Code, Course: r.Course, Quantity: r.Quantity, Status: r.Status, Version: r.Version}
}
func MergeQuantity(a, b Record) Record {
	r := CloneRecord(a)
	if b.Version > a.Version {
		r.Quantity = b.Quantity
		r.Version = b.Version
	}
	return r
}
func Equal(a, b Record) bool {
	return a.ID == b.ID && a.Code == b.Code && a.Course == b.Course && a.Quantity == b.Quantity && a.Status == b.Status && a.Version == b.Version
}
func EnsureStatus(s string) Status {
	switch strings.ToLower(s) {
	case "pending":
		return Pending
	case "approved":
		return Approved
	case "published":
		return Published
	case "archived":
		return Archived
	default:
		return Pending
	}
}
func StatusOrder(s Status) int {
	switch s {
	case Pending:
		return 1
	case Approved:
		return 2
	case Published:
		return 3
	case Archived:
		return 4
	}
	return 0
}
func CompareStatus(a, b Status) int {
	aa, bb := StatusOrder(a), StatusOrder(b)
	if aa < bb {
		return -1
	}
	if aa > bb {
		return 1
	}
	return 0
}
func ValidTransition(r Record, to Status) bool { return CanTransition(r.Status, to) }
func RequireMutable(r Record) error {
	if !IsMutable(r.Status) {
		return errors.New("immutable")
	}
	return nil
}
func CodeKey(r Record) string { return strings.ToUpper(r.Course) + ":" + strings.ToUpper(r.Code) }
func NonEmptyFields(r Record) int {
	n := 0
	for _, v := range []string{r.ID, r.Code, r.Course, string(r.Status)} {
		if strings.TrimSpace(v) != "" {
			n++
		}
	}
	return n
}
func NormalizeRecord(r Record) Record {
	r.Code = NormalizeCode(r.Code)
	r.Course = strings.TrimSpace(r.Course)
	return r
}
func QuantityDelta(a, b Record) int { return b.Quantity - a.Quantity }
func IsZero(r Record) bool {
	return r.ID == "" && r.Code == "" && r.Course == "" && r.Quantity == 0 && r.Version == 0
}
func WithStatus(r Record, s Status) Record { r.Status = s; return r }
func WithVersion(r Record, v int) Record   { r.Version = v; return r }
func NextVersion(v int) int {
	if v < 1 {
		return 1
	}
	return v + 1
}
func Statuses() []Status { return []Status{Pending, Approved, Published, Archived} }
func AllowedStatus(s Status) bool {
	for _, v := range Statuses() {
		if s == v {
			return true
		}
	}
	return false
}
func Describe(r Record) string    { return r.ID + " " + r.Code + " " + string(r.Status) }
func QuantityValid(q int) bool    { return q >= 0 }
func SameCourse(a, b Record) bool { return a.Course == b.Course }
func SameCode(a, b Record) bool   { return NormalizeCode(a.Code) == NormalizeCode(b.Code) }
func RecordScore(r Record) int    { return r.Quantity * StatusOrder(r.Status) }
func PickPreferred(a, b Record) Record {
	if RecordScore(a) >= RecordScore(b) {
		return a
	}
	return b
}
func IsPublished(r Record) bool { return r.Status == Published }
func IsApproved(r Record) bool  { return r.Status == Approved }
func IsArchived(r Record) bool  { return r.Status == Archived }
func IsPending(r Record) bool   { return r.Status == Pending }
func RecordLabels(r Record) []string {
	return []string{r.ID, r.Code, r.Course, StatusName(r.Status), QuantityBand(r.Quantity)}
}
func ValidForSearch(r Record) bool {
	return ValidateID(r.ID) && ValidateCode(r.Code) && ValidateCourse(r.Course)
}
func ClampQuantity(q, max int) int {
	if q < 0 {
		return 0
	}
	if max > 0 && q > max {
		return max
	}
	return q
}
func Increment(r *Record) { r.Quantity++; r.Version = NextVersion(r.Version) }
func Decrement(r *Record) {
	if r.Quantity > 0 {
		r.Quantity--
	}
	r.Version = NextVersion(r.Version)
}
func ResetQuantity(r *Record)           { r.Quantity = 0; r.Version = NextVersion(r.Version) }
func CopyStatuses(in []Status) []Status { return append([]Status(nil), in...) }
func ContainsStatus(in []Status, s Status) bool {
	for _, v := range in {
		if v == s {
			return true
		}
	}
	return false
}
func FilterValid(in []Record) []Record {
	out := []Record{}
	for _, r := range in {
		if r.Valid() {
			out = append(out, r)
		}
	}
	return out
}
func SumVersions(in []Record) int {
	n := 0
	for _, r := range in {
		n += r.Version
	}
	return n
}
func AverageQuantity(in []Record) float64 {
	if len(in) == 0 {
		return 0
	}
	n := 0
	for _, r := range in {
		n += r.Quantity
	}
	return float64(n) / float64(len(in))
}
func SortByStatus(in []Record) []Record {
	out := append([]Record(nil), in...)
	sort.Slice(out, func(i, j int) bool { return StatusOrder(out[i].Status) < StatusOrder(out[j].Status) })
	return out
}
