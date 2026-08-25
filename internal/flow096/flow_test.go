package flow096

import (
	"context"
	"coursecodes/internal/store"
	"path/filepath"
	"testing"
)

func TestWorkflowCreateReviewArchive(t *testing.T) {
	s, _ := store.Open(filepath.Join(t.TempDir(), "d"))
	defer s.Close()
	f := New(s)
	if _, e := f.Register("1", "C1", "Go", 3); e != nil {
		t.Fatal(e)
	}
	if e := f.Workflow(context.Background(), "1"); e != nil {
		t.Fatal(e)
	}
}
func TestWorkflowSearchUpdatePublish(t *testing.T) {
	s, _ := store.Open(filepath.Join(t.TempDir(), "d"))
	defer s.Close()
	f := New(s)
	f.Register("1", "C1", "Go", 3)
	f.Review("1", true)
	if f.Change("1", 5) != nil {
		t.Fatal("change")
	}
	if f.Publish("1") != nil {
		t.Fatal("publish")
	}
}
