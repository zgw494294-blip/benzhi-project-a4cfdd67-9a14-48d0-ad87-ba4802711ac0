package model

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"sort"
)

func HashJSON(v any) string {
	b, _ := json.Marshal(v)
	h := sha256.Sum256(b)
	return hex.EncodeToString(h[:])
}
func ManifestHash(batch DatasetBatch, sources []SourceRecord) string {
	type item struct{ ID, Origin, Checksum, Consent, Sensitivity string }
	ordered := append([]SourceRecord(nil), sources...)
	sort.Slice(ordered, func(i, j int) bool { return ordered[i].SourceID < ordered[j].SourceID })
	items := make([]item, 0, len(ordered))
	for _, s := range ordered {
		items = append(items, item{s.SourceID, s.Origin, s.ContentChecksum, s.ConsentRef, s.Sensitivity})
	}
	return HashJSON(struct {
		BatchID, Policy string
		Scope           []string
		Items           []item
	}{batch.BatchID, batch.PolicyVersion, batch.ReleaseScope, items})
}
