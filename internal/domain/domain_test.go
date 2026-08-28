package domain

import "testing"

func TestRecordLifecycle(t *testing.T) {
	r, e := NewRecord("1", "x", "Go", 2)
	if e != nil {
		t.Fatal(e)
	}
	if r.Approve() != nil || r.Publish() != nil || r.Archive() != nil {
		t.Fatal("lifecycle")
	}
	if !IsTerminal(r.Status) {
		t.Fatal("terminal")
	}
}
