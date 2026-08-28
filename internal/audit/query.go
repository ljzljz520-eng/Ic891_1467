package audit

import "coursecodes/internal/domain"

func FilterByActor(es []domain.AuditEvent, actor string) []domain.AuditEvent {
	out := []domain.AuditEvent{}
	for _, e := range es {
		if e.Actor == actor {
			out = append(out, e)
		}
	}
	return out
}
func Actions(es []domain.AuditEvent) []string {
	out := []string{}
	seen := map[string]bool{}
	for _, e := range es {
		if !seen[e.Action] {
			out = append(out, e.Action)
			seen[e.Action] = true
		}
	}
	return out
}
func Last(es []domain.AuditEvent) (domain.AuditEvent, bool) {
	if len(es) == 0 {
		return domain.AuditEvent{}, false
	}
	return es[len(es)-1], true
}
