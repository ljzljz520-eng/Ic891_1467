package search

import (
	"context"
	"coursecodes/internal/domain"
	"coursecodes/internal/store"
	"path/filepath"
	"testing"
)

func TestFindByCourse(t *testing.T) {
	s, _ := store.Open(filepath.Join(t.TempDir(), "d"))
	defer s.Close()
	r, _ := domain.NewRecord("1", "X", "Go", 1)
	s.SaveRecord(r)
	rows, e := New(s).Find(context.Background(), "go")
	if e != nil || len(rows) != 1 {
		t.Fatal(e, len(rows))
	}
}
