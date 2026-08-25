package flow096

import (
	"context"
	"coursecodes/internal/domain"
)

func (f *Service) QueryAndChange(ctx context.Context, ids []string, quantity int) ([]domain.Record, error) {
	out := []domain.Record{}
	for _, id := range ids {
		if err := ctx.Err(); err != nil {
			return out, err
		}
		r, e := f.Store.GetRecord(id)
		if e != nil {
			return out, e
		}
		// Operate on the record fetched for this id. Reusing a single
		// "selected" variable across iterations caused only the first
		// record to be changed and saved, while later records were left
		// untouched and the first record's quantity was overwritten.
		if e = r.ChangeQuantity(quantity); e != nil {
			return out, e
		}
		if e = f.Store.SaveRecord(r); e != nil {
			return out, e
		}
		out = append(out, r)
	}
	return out, nil
}
