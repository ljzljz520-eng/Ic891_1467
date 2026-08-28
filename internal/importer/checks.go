package importer

import (
	"coursecodes/internal/domain"
	"fmt"
	"strconv"
	"strings"
)

type Validator struct{ Policy domain.Policy }

func NewValidator(p domain.Policy) *Validator { return &Validator{Policy: p} }
func (v *Validator) Check(r domain.Record) error {
	if !r.Valid() {
		return fmt.Errorf("invalid record %s", r.ID)
	}
	return v.Policy.Check(r)
}
func (v *Validator) CheckAll(rows []domain.Record) []string {
	e := []string{}
	for _, r := range rows {
		if x := v.Check(r); x != nil {
			e = append(e, x.Error())
		}
	}
	return e
}
func ParseHeader(line string) []string {
	p := strings.Split(line, ",")
	out := []string{}
	for _, v := range p {
		out = append(out, strings.TrimSpace(strings.ToLower(v)))
	}
	return out
}
func HasHeader(line string) bool {
	return strings.Contains(strings.ToLower(line), "code") && strings.Contains(strings.ToLower(line), "course")
}
func CanonicalRow(r Row) Row {
	return Row{ID: strings.TrimSpace(r.ID), Code: strings.ToUpper(strings.TrimSpace(r.Code)), Course: strings.TrimSpace(r.Course), Quantity: strings.TrimSpace(r.Quantity)}
}
func CanonicalRows(rows []Row) []Row {
	out := []Row{}
	for _, r := range rows {
		out = append(out, CanonicalRow(r))
	}
	return out
}
func RowKey(r Row) string { return r.ID + ":" + strings.ToUpper(r.Code) }
func UniqueRows(rows []Row) []Row {
	seen := map[string]bool{}
	out := []Row{}
	for _, r := range rows {
		if !seen[RowKey(r)] {
			seen[RowKey(r)] = true
			out = append(out, r)
		}
	}
	return out
}
func RowsForCourse(rows []Row, c string) []Row {
	out := []Row{}
	for _, r := range rows {
		if strings.EqualFold(r.Course, c) {
			out = append(out, r)
		}
	}
	return out
}
func RowsForCode(rows []Row, c string) []Row {
	out := []Row{}
	for _, r := range rows {
		if strings.EqualFold(r.Code, c) {
			out = append(out, r)
		}
	}
	return out
}
func CountRows(rows []Row) int { return len(rows) }
func ValidRows(rows []Row) []Row {
	out := []Row{}
	for _, r := range rows {
		if r.ID != "" && r.Code != "" && r.Course != "" {
			out = append(out, r)
		}
	}
	return out
}
func InvalidRows(rows []Row) []Row {
	out := []Row{}
	for _, r := range rows {
		if r.ID == "" || r.Code == "" || r.Course == "" {
			out = append(out, r)
		}
	}
	return out
}
func QuantityRows(rows []Row) int {
	n := 0
	for _, r := range rows {
		q, e := strconv.Atoi(r.Quantity)
		if e == nil {
			n += q
		}
	}
	return n
}
func FormatRows(rows []Row) string {
	p := []string{}
	for _, r := range rows {
		p = append(p, RowKey(r))
	}
	return strings.Join(p, ";")
}
func ErrorReport(errs []string) string {
	if len(errs) == 0 {
		return "ok"
	}
	return strings.Join(errs, ";")
}
func Importable(lines []string) bool { return len(Parse(lines)) > 0 }
func LineCount(lines []string) int   { return len(NormalizeLines(lines)) }
func FirstRow(lines []string) (Row, bool) {
	r := Parse(lines)
	if len(r) == 0 {
		return Row{}, false
	}
	return r[0], true
}
func LastRow(lines []string) (Row, bool) {
	r := Parse(lines)
	if len(r) == 0 {
		return Row{}, false
	}
	return r[len(r)-1], true
}
func ContainsCourse(rows []Row, c string) bool { return len(RowsForCourse(rows, c)) > 0 }
func ContainsCode(rows []Row, c string) bool   { return len(RowsForCode(rows, c)) > 0 }
func ReportLabel(r Report) string              { return FormatReport(r) }
