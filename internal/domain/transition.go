package domain

type Transition struct {
	From, To Status
	Reason   string
}

func (r Record) Transitions() []Transition {
	all := []Transition{{Pending, Approved, "review"}, {Approved, Published, "publish"}, {Published, Archived, "archive"}}
	out := []Transition{}
	for _, t := range all {
		if t.From == r.Status {
			out = append(out, t)
		}
	}
	return out
}
func IsTerminal(s Status) bool { return s == Archived }
func CanTransition(from, to Status) bool {
	switch from {
	case Pending:
		return to == Approved
	case Approved:
		return to == Published
	case Published:
		return to == Archived
	default:
		return false
	}
}
