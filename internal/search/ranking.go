package search

import (
	"coursecodes/internal/domain"
	"sort"
	"strings"
)

type Result struct {
	Record domain.Record
	Score  int
}

func Rank(rows []domain.Record, q string) []Result {
	out := []Result{}
	for _, r := range rows {
		score := 0
		if strings.EqualFold(r.Code, q) {
			score += 100
		}
		if strings.Contains(strings.ToLower(r.Course), strings.ToLower(q)) {
			score += 20
		}
		score += r.Quantity
		out = append(out, Result{r, score})
	}
	sort.Slice(out, func(i, j int) bool { return out[i].Score > out[j].Score })
	return out
}
func Top(rows []domain.Record, q string, n int) []domain.Record {
	rank := Rank(rows, q)
	if n > len(rank) {
		n = len(rank)
	}
	out := []domain.Record{}
	for _, v := range rank[:n] {
		out = append(out, v.Record)
	}
	return out
}
func QueryCodes(rows []domain.Record) map[string]domain.Record {
	m := map[string]domain.Record{}
	for _, r := range rows {
		m[strings.ToUpper(r.Code)] = r
	}
	return m
}
func DistinctCourses(rows []domain.Record) []string {
	m := map[string]bool{}
	out := []string{}
	for _, r := range rows {
		if !m[r.Course] {
			m[r.Course] = true
			out = append(out, r.Course)
		}
	}
	sort.Strings(out)
	return out
}
func DistinctCodes(rows []domain.Record) []string {
	m := map[string]bool{}
	out := []string{}
	for _, r := range rows {
		c := strings.ToUpper(r.Code)
		if !m[c] {
			m[c] = true
			out = append(out, c)
		}
	}
	sort.Strings(out)
	return out
}
func Prefix(rows []domain.Record, p string) []domain.Record {
	out := []domain.Record{}
	for _, r := range rows {
		if strings.HasPrefix(strings.ToUpper(r.Code), strings.ToUpper(p)) {
			out = append(out, r)
		}
	}
	return out
}
func Suffix(rows []domain.Record, p string) []domain.Record {
	out := []domain.Record{}
	for _, r := range rows {
		if strings.HasSuffix(strings.ToUpper(r.Code), strings.ToUpper(p)) {
			out = append(out, r)
		}
	}
	return out
}
func Exact(rows []domain.Record, c string) []domain.Record {
	out := []domain.Record{}
	for _, r := range rows {
		if strings.EqualFold(r.Code, c) {
			out = append(out, r)
		}
	}
	return out
}
func QuantityAtLeast(rows []domain.Record, n int) []domain.Record {
	out := []domain.Record{}
	for _, r := range rows {
		if r.Quantity >= n {
			out = append(out, r)
		}
	}
	return out
}
func QuantityAtMost(rows []domain.Record, n int) []domain.Record {
	out := []domain.Record{}
	for _, r := range rows {
		if r.Quantity <= n {
			out = append(out, r)
		}
	}
	return out
}
func SearchTerms(rows []domain.Record, terms []string) []domain.Record {
	out := []domain.Record{}
	for _, r := range rows {
		ok := true
		text := strings.ToLower(r.Code + " " + r.Course)
		for _, t := range terms {
			if !strings.Contains(text, strings.ToLower(t)) {
				ok = false
			}
		}
		if ok {
			out = append(out, r)
		}
	}
	return out
}
func Paginate(rows []domain.Record, page, size int) []domain.Record {
	if page < 1 {
		page = 1
	}
	if size < 1 {
		size = 1
	}
	start := (page - 1) * size
	if start >= len(rows) {
		return []domain.Record{}
	}
	end := start + size
	if end > len(rows) {
		end = len(rows)
	}
	return rows[start:end]
}
func PageCount(total, size int) int {
	if size <= 0 {
		return 0
	}
	return (total + size - 1) / size
}
func SortCourse(rows []domain.Record) []domain.Record {
	out := append([]domain.Record(nil), rows...)
	sort.Slice(out, func(i, j int) bool { return out[i].Course < out[j].Course })
	return out
}
func SortCode(rows []domain.Record) []domain.Record {
	out := append([]domain.Record(nil), rows...)
	sort.Slice(out, func(i, j int) bool { return out[i].Code < out[j].Code })
	return out
}
func SortNewest(rows []domain.Record) []domain.Record {
	out := append([]domain.Record(nil), rows...)
	sort.Slice(out, func(i, j int) bool { return out[i].Version > out[j].Version })
	return out
}
func Score(r domain.Record, q string) int { return Rank([]domain.Record{r}, q)[0].Score }
func Best(rows []domain.Record, q string) (domain.Record, bool) {
	x := Top(rows, q, 1)
	if len(x) == 0 {
		return domain.Record{}, false
	}
	return x[0], true
}
func Any(rows []domain.Record, q string) bool {
	for _, r := range rows {
		if Match(r, q) {
			return true
		}
	}
	return false
}
func Empty(rows []domain.Record) bool { return len(rows) == 0 }
func Count(rows []domain.Record) int  { return len(rows) }
