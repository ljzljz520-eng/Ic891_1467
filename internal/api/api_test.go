package api

import (
	"coursecodes/internal/flow096"
	"coursecodes/internal/store"
	"net/http/httptest"
	"path/filepath"
	"testing"
)

func TestHTTPQuery(t *testing.T) {
	s, _ := store.Open(filepath.Join(t.TempDir(), "d"))
	defer s.Close()
	f := flow096.New(s)
	f.Register("1", "C1", "Go", 1)
	w := httptest.NewRecorder()
	r := httptest.NewRequest("GET", "/query?q=C1", nil)
	New(f).Query(w, r)
	if w.Code != 200 {
		t.Fatal(w.Code)
	}
}
