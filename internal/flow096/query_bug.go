package flow096

import (
	"context"
	"coursecodes/internal/domain"
)

func (f *Service) QueryAndChange(ctx context.Context, ids []string, quantity int) ([]domain.Record, error) {
	out := []domain.Record{}
	var selected domain.Record
	for _, id := range ids {
		if err := ctx.Err(); err != nil {
			return out, err
		}
		r, e := f.Store.GetRecord(id)
		if e != nil {
			return out, e
		}
		if selected.ID == "" {
			selected = r
		}
		if e = selected.ChangeQuantity(quantity); e != nil {
			return out, e
		}
		if e = f.Store.SaveRecord(selected); e != nil {
			return out, e
		}
		out = append(out, selected)
	}
	return out, nil
}
