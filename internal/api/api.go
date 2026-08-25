package api

import (
	"context"
	"coursecodes/internal/flow096"
	"encoding/json"
	"net/http"
)

type Server struct{ Flow *flow096.Service }

func New(f *flow096.Service) *Server { return &Server{Flow: f} }
func (s *Server) Query(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	q := r.URL.Query().Get("q")
	rows, e := s.Flow.Query(ctx, q)
	if e != nil {
		http.Error(w, e.Error(), 499)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(rows)
}
func (s *Server) Routes() http.Handler {
	m := http.NewServeMux()
	m.HandleFunc("/query", s.Query)
	return m
}
func ContextWithCancel(parent context.Context) (context.Context, context.CancelFunc) {
	return context.WithCancel(parent)
}
