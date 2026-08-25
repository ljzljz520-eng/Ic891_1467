package audit

import (
	"coursecodes/internal/store"
	"path/filepath"
	"testing"
)

func TestAuditSummary(t *testing.T) {
	s, _ := store.Open(filepath.Join(t.TempDir(), "d"))
	defer s.Close()
	l := New(s)
	l.Record("1", "register", "a")
	l.Record("1", "publish", "b")
	es, _ := l.ForRecord("1")
	if !HasAction(es, "publish") || Summarize(es)["register"] != 1 {
		t.Fatal(es)
	}
}
