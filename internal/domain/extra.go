package domain

func (r Record) CanReview() bool        { return r.Status == Pending }
func (r Record) CanPublish() bool       { return r.Status == Approved }
func (r Record) CanArchive() bool       { return r.Status == Published }
func (r Record) IsActive() bool         { return r.Status != Archived }
func (r Record) IsDraft() bool          { return r.Status == Pending }
func (r Record) IsReady() bool          { return r.Status == Approved || r.Status == Published }
func (r Record) QuantityPositive() bool { return r.Quantity > 0 }
func (r Record) QuantityZero() bool     { return r.Quantity == 0 }
func (r Record) VersionPositive() bool  { return r.Version > 0 }
func (r Record) CourseKey() string      { return r.Course }
func (r Record) CodeKey() string        { return NormalizeCode(r.Code) }
func (r Record) Identity() string       { return r.ID }
func (r Record) Summary() string        { return Describe(r) }
func (r Record) Labels() []string       { return RecordLabels(r) }
func (r Record) Score() int             { return RecordScore(r) }
func (r Record) NextStatus() Status {
	for _, t := range r.Transitions() {
		return t.To
	}
	return r.Status
}
func (r Record) TransitionReason() string {
	for _, t := range r.Transitions() {
		return t.Reason
	}
	return ""
}
func (r Record) FieldCount() int               { return NonEmptyFields(r) }
func (r Record) PolicyValid() bool             { return DefaultPolicy().Check(r) == nil }
func (r Record) Searchable() bool              { return ValidForSearch(r) }
func (r Record) QuantityBand() string          { return QuantityBand(r.Quantity) }
func (r Record) StatusRank() int               { return StatusOrder(r.Status) }
func (r Record) IsTerminal() bool              { return IsTerminal(r.Status) }
func (r Record) CanChange() bool               { return IsMutable(r.Status) }
func (r Record) Clone() Record                 { return CloneRecord(r) }
func (r Record) WithQuantity(q int) Record     { r.Quantity = q; return r }
func (r Record) WithCode(c string) Record      { r.Code = c; return r }
func (r Record) WithCourse(c string) Record    { r.Course = c; return r }
func (r Record) Bump() Record                  { r.Version = NextVersion(r.Version); return r }
func (r Record) Normalize() Record             { return NormalizeRecord(r) }
func (r Record) ValidID() bool                 { return ValidateID(r.ID) }
func (r Record) ValidCode() bool               { return ValidateCode(r.Code) }
func (r Record) ValidCourse() bool             { return ValidateCourse(r.Course) }
func (r Record) ValidQuantity() bool           { return QuantityValid(r.Quantity) }
func (r Record) SameCode(other Record) bool    { return SameCode(r, other) }
func (r Record) SameCourse(other Record) bool  { return SameCourse(r, other) }
func (r Record) Delta(other Record) int        { return QuantityDelta(r, other) }
func (r Record) Preferred(other Record) Record { return PickPreferred(r, other) }
func (r Record) Mutate(q int) error            { return r.ChangeQuantity(q) }
func (r Record) StatusString() string          { return StatusName(r.Status) }
func (r Record) IsStatus(s Status) bool        { return r.Status == s }
func RecordsValid(rs []Record) bool {
	for _, r := range rs {
		if !r.Valid() {
			return false
		}
	}
	return true
}
func RecordsTerminal(rs []Record) bool {
	for _, r := range rs {
		if !r.IsTerminal() {
			return false
		}
	}
	return true
}
func RecordsActive(rs []Record) int {
	n := 0
	for _, r := range rs {
		if r.IsActive() {
			n++
		}
	}
	return n
}
func RecordsQuantity(rs []Record) int {
	n := 0
	for _, r := range rs {
		n += r.Quantity
	}
	return n
}
func RecordsVersions(rs []Record) int {
	n := 0
	for _, r := range rs {
		n += r.Version
	}
	return n
}
func RecordsByStatus(rs []Record, s Status) []Record {
	out := []Record{}
	for _, r := range rs {
		if r.Status == s {
			out = append(out, r)
		}
	}
	return out
}
func RecordsByCourse(rs []Record, c string) []Record {
	out := []Record{}
	for _, r := range rs {
		if r.Course == c {
			out = append(out, r)
		}
	}
	return out
}
func RecordsByCode(rs []Record, c string) []Record {
	out := []Record{}
	for _, r := range rs {
		if SameCode(r, Record{Code: c}) {
			out = append(out, r)
		}
	}
	return out
}
func RecordAt(rs []Record, i int) (Record, bool) {
	if i < 0 || i >= len(rs) {
		return Record{}, false
	}
	return rs[i], true
}
func FirstRecord(rs []Record) (Record, bool) { return RecordAt(rs, 0) }
func LastRecord(rs []Record) (Record, bool)  { return RecordAt(rs, len(rs)-1) }
func RecordCount(rs []Record) int            { return len(rs) }
