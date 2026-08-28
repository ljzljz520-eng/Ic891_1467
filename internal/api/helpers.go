package api

import (
	"coursecodes/internal/domain"
	"encoding/json"
	"net/http"
	"strings"
)

type Envelope struct {
	OK    bool
	Error string
	Data  any
}

func WriteJSON(w http.ResponseWriter, v any) {
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(v)
}
func WriteError(w http.ResponseWriter, e error) {
	w.WriteHeader(http.StatusBadRequest)
	WriteJSON(w, Envelope{OK: false, Error: e.Error()})
}
func MethodAllowed(r *http.Request, method string) bool { return r.Method == method }
func ParseID(r *http.Request) string                    { return strings.TrimSpace(r.URL.Query().Get("id")) }
func ParseQuantity(r *http.Request) int {
	var v struct{ Quantity int }
	_ = json.NewDecoder(r.Body).Decode(&v)
	return v.Quantity
}
func RespondRecords(w http.ResponseWriter, rs []domain.Record) {
	WriteJSON(w, Envelope{OK: true, Data: rs})
}
func RespondRecord(w http.ResponseWriter, r domain.Record) { WriteJSON(w, Envelope{OK: true, Data: r}) }
func Header(w http.ResponseWriter, k, v string)            { w.Header().Set(k, v) }
func CORS(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		Header(w, "Access-Control-Allow-Origin", "*")
		next.ServeHTTP(w, r)
	})
}
func NoCache(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		Header(w, "Cache-Control", "no-store")
		next.ServeHTTP(w, r)
	})
}
func Chain(h http.Handler, ms ...func(http.Handler) http.Handler) http.Handler {
	for i := len(ms) - 1; i >= 0; i-- {
		h = ms[i](h)
	}
	return h
}
func StatusCode(e error) int {
	if e == nil {
		return http.StatusOK
	}
	return http.StatusBadRequest
}
func IsJSON(r *http.Request) bool {
	return strings.Contains(r.Header.Get("Content-Type"), "application/json")
}
func QueryValue(r *http.Request, k string) string { return strings.TrimSpace(r.URL.Query().Get(k)) }
func BoolValue(r *http.Request, k string) bool    { return QueryValue(r, k) == "true" }
func Empty(v string) bool                         { return strings.TrimSpace(v) == "" }
func ValidID(v string) bool                       { return !Empty(v) && len(v) < 100 }
func ValidQuantity(q int) bool                    { return q >= 0 }
func ValidStatus(s domain.Status) bool            { return domain.AllowedStatus(s) }
func RecordSummary(r domain.Record) map[string]any {
	return map[string]any{"id": r.ID, "code": r.Code, "quantity": r.Quantity, "status": r.Status}
}
func RecordsSummary(rs []domain.Record) []map[string]any {
	out := []map[string]any{}
	for _, r := range rs {
		out = append(out, RecordSummary(r))
	}
	return out
}
func ParseRecord(data []byte) (domain.Record, error) {
	var r domain.Record
	e := json.Unmarshal(data, &r)
	return r, e
}
func EncodeRecord(r domain.Record) []byte { b, _ := json.Marshal(r); return b }
func RouteName(path string) string {
	p := strings.Trim(path, "/")
	if p == "" {
		return "root"
	}
	return p
}
func NormalizePath(path string) string { return "/" + strings.Trim(strings.TrimSpace(path), "/") }
func HasPrefix(path, prefix string) bool {
	return strings.HasPrefix(NormalizePath(path), NormalizePath(prefix))
}
func SplitPath(path string) []string {
	p := strings.Split(strings.Trim(path, "/"), "/")
	if len(p) == 1 && p[0] == "" {
		return []string{}
	}
	return p
}
func LastPath(path string) string {
	p := SplitPath(path)
	if len(p) == 0 {
		return ""
	}
	return p[len(p)-1]
}
func PathDepth(path string) int              { return len(SplitPath(path)) }
func Accepts(r *http.Request, v string) bool { return strings.Contains(r.Header.Get("Accept"), v) }
func UserAgent(r *http.Request) string       { return r.Header.Get("User-Agent") }
func RequestID(r *http.Request) string       { return r.Header.Get("X-Request-ID") }
func WithDefaultRequestID(r *http.Request) string {
	if id := RequestID(r); id != "" {
		return id
	}
	return "request-static"
}
func IsCanceled(r *http.Request) bool     { return r.Context().Err() != nil }
func ContentLength(r *http.Request) int64 { return r.ContentLength }
func RecordStatus(r domain.Record) string { return domain.StatusName(r.Status) }
func RecordCode(r domain.Record) string   { return strings.ToUpper(r.Code) }
func RecordCourse(r domain.Record) string { return strings.TrimSpace(r.Course) }
func RecordQuantity(r domain.Record) int  { return r.Quantity }
func RecordVersion(r domain.Record) int   { return r.Version }
func IsArchived(r domain.Record) bool     { return r.Status == domain.Archived }
func IsActive(r domain.Record) bool       { return r.Status != domain.Archived }
func FilterActive(rs []domain.Record) []domain.Record {
	out := []domain.Record{}
	for _, r := range rs {
		if IsActive(r) {
			out = append(out, r)
		}
	}
	return out
}
func FilterArchived(rs []domain.Record) []domain.Record {
	out := []domain.Record{}
	for _, r := range rs {
		if IsArchived(r) {
			out = append(out, r)
		}
	}
	return out
}
func QuantityTotal(rs []domain.Record) int {
	n := 0
	for _, r := range rs {
		n += r.Quantity
	}
	return n
}
func QuantityMaximum(rs []domain.Record) int {
	n := 0
	for _, r := range rs {
		if r.Quantity > n {
			n = r.Quantity
		}
	}
	return n
}
func QuantityMinimum(rs []domain.Record) int {
	if len(rs) == 0 {
		return 0
	}
	n := rs[0].Quantity
	for _, r := range rs[1:] {
		if r.Quantity < n {
			n = r.Quantity
		}
	}
	return n
}
func RecordCount(rs []domain.Record) int { return len(rs) }
func ContainsRecord(rs []domain.Record, id string) bool {
	for _, r := range rs {
		if r.ID == id {
			return true
		}
	}
	return false
}
func FindRecord(rs []domain.Record, id string) (domain.Record, bool) {
	for _, r := range rs {
		if r.ID == id {
			return r, true
		}
	}
	return domain.Record{}, false
}
func UpdateLocal(rs []domain.Record, id string, q int) []domain.Record {
	out := append([]domain.Record(nil), rs...)
	for i := range out {
		if out[i].ID == id {
			out[i].Quantity = q
			out[i].Version++
		}
	}
	return out
}
func AddRecord(rs []domain.Record, r domain.Record) []domain.Record {
	return append(append([]domain.Record(nil), rs...), r)
}
func RemoveRecord(rs []domain.Record, id string) []domain.Record {
	out := []domain.Record{}
	for _, r := range rs {
		if r.ID != id {
			out = append(out, r)
		}
	}
	return out
}
func CloneRecords(rs []domain.Record) []domain.Record { return append([]domain.Record(nil), rs...) }
func SortRecords(rs []domain.Record) []domain.Record {
	out := CloneRecords(rs)
	for i := 0; i < len(out); i++ {
		for j := i + 1; j < len(out); j++ {
			if out[j].ID < out[i].ID {
				out[i], out[j] = out[j], out[i]
			}
		}
	}
	return out
}
