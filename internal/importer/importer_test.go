package importer

import (
	"coursecodes/internal/store"
	"path/filepath"
	"testing"
)

func TestImportReport(t *testing.T) {
	s, _ := store.Open(filepath.Join(t.TempDir(), "d"))
	defer s.Close()
	r, e := Import(s, []string{"1,C1,Go,3", "2,C2,SQL,no"})
	if e != nil || r.Imported != 1 || r.Rejected != 1 {
		t.Fatal(r, e)
	}
}
