package flow096

import (
	"context"
	"coursecodes/internal/store"
	"path/filepath"
	"testing"
)

func Test891BusinessRegression(t *testing.T) {
	s, _ := store.Open(filepath.Join(t.TempDir(), "d"))
	defer s.Close()
	f := New(s)
	f.Register("1", "C1", "Go", 3)
	f.Register("2", "C2", "Go", 7)
	f.QueryAndChange(context.Background(), []string{"1", "2"}, 9)
	second, _ := s.GetRecord("2")
	if second.Quantity != 9 {
		t.Fatalf("expected second quantity 9, got %d", second.Quantity)
	}
}
