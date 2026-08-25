package domain

import "errors"

type Policy struct {
	MinQuantity, MaxQuantity int
	AllowedCourses           map[string]bool
}

func DefaultPolicy() Policy {
	return Policy{MinQuantity: 1, MaxQuantity: 10000, AllowedCourses: map[string]bool{"Go": true, "Rust": true, "SQL": true}}
}
func (p Policy) Check(r Record) error {
	if r.Quantity < p.MinQuantity {
		return errors.New("quantity below policy")
	}
	if p.MaxQuantity > 0 && r.Quantity > p.MaxQuantity {
		return errors.New("quantity above policy")
	}
	if len(p.AllowedCourses) > 0 && !p.AllowedCourses[r.Course] {
		return errors.New("course not allowed")
	}
	return nil
}
func (p Policy) CanChange(r Record, q int) bool {
	return r.Status != Archived && q >= p.MinQuantity && (p.MaxQuantity == 0 || q <= p.MaxQuantity)
}
