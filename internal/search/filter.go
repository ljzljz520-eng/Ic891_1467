package search

import "coursecodes/internal/domain"

type Filter struct {
	Course, Code string
	Status       domain.Status
	MinQuantity  int
}

func Apply(rows []domain.Record, f Filter) []domain.Record {
	out := []domain.Record{}
	for _, r := range rows {
		if f.Course != "" && r.Course != f.Course {
			continue
		}
		if f.Code != "" && !Match(r, f.Code) {
			continue
		}
		if f.Status != "" && r.Status != f.Status {
			continue
		}
		if r.Quantity < f.MinQuantity {
			continue
		}
		out = append(out, r)
	}
	return out
}
func GroupByCourse(rows []domain.Record) map[string][]domain.Record {
	m := map[string][]domain.Record{}
	for _, r := range rows {
		m[r.Course] = append(m[r.Course], r)
	}
	return m
}
func HighestQuantity(rows []domain.Record) (domain.Record, bool) {
	var best domain.Record
	ok := false
	for _, r := range rows {
		if !ok || r.Quantity > best.Quantity {
			best = r
			ok = true
		}
	}
	return best, ok
}
