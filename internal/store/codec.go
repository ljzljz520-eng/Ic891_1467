package store

import (
	"coursecodes/internal/domain"
	"encoding/base64"
	"encoding/json"
)

func EncodeRecord(r domain.Record) (string, error) {
	b, e := json.Marshal(r)
	if e != nil {
		return "", e
	}
	return base64.StdEncoding.EncodeToString(b), nil
}
func DecodeRecord(s string) (domain.Record, error) {
	b, e := base64.StdEncoding.DecodeString(s)
	if e != nil {
		return domain.Record{}, e
	}
	var r domain.Record
	e = json.Unmarshal(b, &r)
	return r, e
}
func EncodeRecords(rs []domain.Record) ([]byte, error) { return json.Marshal(rs) }
func DecodeRecords(b []byte) ([]domain.Record, error) {
	var rs []domain.Record
	e := json.Unmarshal(b, &rs)
	return rs, e
}
func EncodeEvent(e domain.AuditEvent) ([]byte, error) { return json.Marshal(e) }
func DecodeEvent(b []byte) (domain.AuditEvent, error) {
	var e domain.AuditEvent
	x := json.Unmarshal(b, &e)
	return e, x
}
func EncodeWorkflow(w domain.Workflow) ([]byte, error) { return json.Marshal(w) }
func DecodeWorkflow(b []byte) (domain.Workflow, error) {
	var w domain.Workflow
	x := json.Unmarshal(b, &w)
	return w, x
}
func EncodeAttachment(a domain.Attachment) ([]byte, error) { return json.Marshal(a) }
func DecodeAttachment(b []byte) (domain.Attachment, error) {
	var a domain.Attachment
	x := json.Unmarshal(b, &a)
	return a, x
}
func RecordBytes(r domain.Record) []byte              { b, _ := json.Marshal(r); return b }
func EventBytes(e domain.AuditEvent) []byte           { b, _ := json.Marshal(e); return b }
func WorkflowBytes(w domain.Workflow) []byte          { b, _ := json.Marshal(w); return b }
func AttachmentBytes(a domain.Attachment) []byte      { b, _ := json.Marshal(a); return b }
func CloneRecords(in []domain.Record) []domain.Record { return append([]domain.Record(nil), in...) }
func ReplaceRecord(in []domain.Record, r domain.Record) []domain.Record {
	out := []domain.Record{}
	found := false
	for _, v := range in {
		if v.ID == r.ID {
			out = append(out, r)
			found = true
		} else {
			out = append(out, v)
		}
	}
	if !found {
		out = append(out, r)
	}
	return out
}
func RemoveRecord(in []domain.Record, id string) []domain.Record {
	out := []domain.Record{}
	for _, r := range in {
		if r.ID != id {
			out = append(out, r)
		}
	}
	return out
}
func FindRecord(in []domain.Record, id string) (domain.Record, bool) {
	for _, r := range in {
		if r.ID == id {
			return r, true
		}
	}
	return domain.Record{}, false
}
func MarshalMap(m map[string]string) []byte { b, _ := json.Marshal(m); return b }
func UnmarshalMap(b []byte) map[string]string {
	var m map[string]string
	_ = json.Unmarshal(b, &m)
	return m
}
func EncodeText(s string) string { return base64.StdEncoding.EncodeToString([]byte(s)) }
func DecodeText(s string) string { b, _ := base64.StdEncoding.DecodeString(s); return string(b) }
func RecordDigest(r domain.Record) string {
	b := RecordBytes(r)
	return base64.RawStdEncoding.EncodeToString(b)
}
func BatchDigest(rs []domain.Record) string {
	b, _ := EncodeRecords(rs)
	return base64.RawStdEncoding.EncodeToString(b)
}
func EnsureRecord(r domain.Record) domain.Record {
	if r.Version == 0 {
		r.Version = 1
	}
	if r.Status == "" {
		r.Status = domain.Pending
	}
	return r
}
