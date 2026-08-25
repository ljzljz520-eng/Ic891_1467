package importer

import (
	"sort"
	"strings"
)

func NormalizeLines(lines []string) []string {
	out := []string{}
	for _, l := range lines {
		l = strings.TrimSpace(l)
		if l != "" {
			out = append(out, l)
		}
	}
	return out
}
func SortErrors(errs []string) []string {
	out := append([]string(nil), errs...)
	sort.Strings(out)
	return out
}
func SuccessRate(r Report) float64 {
	total := r.Imported + r.Rejected
	if total == 0 {
		return 1
	}
	return float64(r.Imported) / float64(total)
}
