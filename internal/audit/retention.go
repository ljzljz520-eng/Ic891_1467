package audit

import (
	"coursecodes/internal/domain"
	"strings"
)

type Timeline struct{ Events []domain.AuditEvent }

func NewTimeline(e []domain.AuditEvent) *Timeline {
	return &Timeline{Events: append([]domain.AuditEvent(nil), e...)}
}
func (t *Timeline) Len() int               { return len(t.Events) }
func (t *Timeline) Empty() bool            { return len(t.Events) == 0 }
func (t *Timeline) Actions() []string      { return Actions(t.Events) }
func (t *Timeline) Contains(a string) bool { return HasAction(t.Events, a) }
func (t *Timeline) Actors() []string {
	out := []string{}
	for _, e := range t.Events {
		found := false
		for _, v := range out {
			if v == e.Actor {
				found = true
			}
		}
		if !found {
			out = append(out, e.Actor)
		}
	}
	return out
}
func (t *Timeline) ByActor(a string) *Timeline { return NewTimeline(FilterByActor(t.Events, a)) }
func (t *Timeline) LastAction() string {
	e, ok := Last(t.Events)
	if !ok {
		return ""
	}
	return e.Action
}
func (t *Timeline) FirstAction() string {
	if len(t.Events) == 0 {
		return ""
	}
	return t.Events[0].Action
}
func (t *Timeline) CountAction(a string) int {
	n := 0
	for _, e := range t.Events {
		if e.Action == a {
			n++
		}
	}
	return n
}
func (t *Timeline) RecordIDs() []string {
	out := []string{}
	for _, e := range t.Events {
		out = append(out, e.RecordID)
	}
	return out
}
func (t *Timeline) FilterRecord(id string) *Timeline {
	out := []domain.AuditEvent{}
	for _, e := range t.Events {
		if e.RecordID == id {
			out = append(out, e)
		}
	}
	return NewTimeline(out)
}
func (t *Timeline) Merge(other *Timeline) *Timeline {
	out := append([]domain.AuditEvent{}, t.Events...)
	out = append(out, other.Events...)
	return NewTimeline(out)
}
func (t *Timeline) Clone() *Timeline { return NewTimeline(t.Events) }
func (t *Timeline) HasActor(a string) bool {
	for _, v := range t.Actors() {
		if v == a {
			return true
		}
	}
	return false
}
func (t *Timeline) HasRecord(id string) bool {
	for _, v := range t.RecordIDs() {
		if v == id {
			return true
		}
	}
	return false
}
func (t *Timeline) Summary() map[string]int { return Summarize(t.Events) }
func (t *Timeline) LatestFor(id string) (domain.AuditEvent, bool) {
	return Last(t.FilterRecord(id).Events)
}
func (t *Timeline) IsComplete() bool   { return t.Contains("archive") }
func (t *Timeline) IsPublished() bool  { return t.Contains("publish") }
func (t *Timeline) IsRegistered() bool { return t.Contains("register") }
func (t *Timeline) IsChanged() bool    { return t.Contains("change") }
func (t *Timeline) IsApproved() bool   { return t.Contains("approve") }
func (t *Timeline) SortByAction() []domain.AuditEvent {
	out := append([]domain.AuditEvent(nil), t.Events...)
	for i := range out {
		for j := i + 1; j < len(out); j++ {
			if out[j].Action < out[i].Action {
				out[i], out[j] = out[j], out[i]
			}
		}
	}
	return out
}
func (t *Timeline) Limit(n int) *Timeline {
	if n < 0 {
		n = 0
	}
	if n > len(t.Events) {
		n = len(t.Events)
	}
	return NewTimeline(t.Events[:n])
}
func (t *Timeline) Tail(n int) *Timeline {
	if n < 0 {
		n = 0
	}
	if n > len(t.Events) {
		n = len(t.Events)
	}
	return NewTimeline(t.Events[len(t.Events)-n:])
}
func (t *Timeline) Text() string {
	p := []string{}
	for _, e := range t.Events {
		p = append(p, e.Action+":"+e.Actor)
	}
	return strings.Join(p, "|")
}
