package importer

import (
	"coursecodes/internal/domain"
	"coursecodes/internal/store"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"strconv"
	"strings"
)

type Row struct{ ID, Code, Course, Quantity string }
type Report struct {
	Imported, Rejected int
	Errors             []string
}

func Parse(lines []string) []Row {
	rows := []Row{}
	for _, line := range lines {
		p := strings.Split(line, ",")
		if len(p) == 4 {
			rows = append(rows, Row{p[0], p[1], p[2], p[3]})
		}
	}
	return rows
}
func Import(s *store.Store, lines []string) (Report, error) {
	r := Report{}
	for i, row := range Parse(lines) {
		q, e := strconv.Atoi(row.Quantity)
		if e != nil || q < 0 {
			r.Rejected++
			r.Errors = append(r.Errors, fmt.Sprintf("row %d quantity", i+1))
			continue
		}
		rec, e := domain.NewRecord(row.ID, domain.NormalizeCode(row.Code), row.Course, q)
		if e != nil {
			r.Rejected++
			r.Errors = append(r.Errors, e.Error())
			continue
		}
		if e = s.SaveRecord(rec); e != nil {
			return r, e
		}
		d := sha256.Sum256([]byte(row.ID + row.Code))
		a := domain.Attachment{ID: "attachment-" + row.ID, RecordID: row.ID, Name: "import", Digest: hex.EncodeToString(d[:])}
		if e = s.SaveAttachment(a); e != nil {
			return r, e
		}
		r.Imported++
	}
	return r, nil
}
func Validate(rows []Row) []string {
	errs := []string{}
	seen := map[string]bool{}
	for i, r := range rows {
		if seen[r.Code] {
			errs = append(errs, fmt.Sprintf("row %d duplicate", i+1))
		}
		seen[r.Code] = true
	}
	return errs
}
func FormatReport(r Report) string {
	return fmt.Sprintf("imported=%d rejected=%d", r.Imported, r.Rejected)
}
