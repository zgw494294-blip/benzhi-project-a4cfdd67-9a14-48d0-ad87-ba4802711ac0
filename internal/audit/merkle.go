package audit

import "github.com/benzhi-project/a4cfdd67-9a14-48d0-ad87-ba4802711ac0/internal/model"

func RootHash(events []model.AuditEvent) string {
	if len(events) == 0 {
		return ""
	}
	level := make([]string, 0, len(events))
	for _, e := range events {
		level = append(level, e.EventHash)
	}
	for len(level) > 1 {
		next := []string{}
		for i := 0; i < len(level); i += 2 {
			j := i + 1
			if j >= len(level) {
				j = i
			}
			next = append(next, model.HashJSON([]string{level[i], level[j]}))
		}
		level = next
	}
	return level[0]
}
func InclusionProof(events []model.AuditEvent, index int) []string {
	if index < 0 || index >= len(events) {
		return nil
	}
	level := make([]string, 0, len(events))
	for _, e := range events {
		level = append(level, e.EventHash)
	}
	proof := []string{}
	pos := index
	for len(level) > 1 {
		if pos%2 == 0 {
			if pos+1 < len(level) {
				proof = append(proof, level[pos+1])
			} else {
				proof = append(proof, level[pos])
			}
		} else {
			proof = append(proof, level[pos-1])
		}
		next := []string{}
		for i := 0; i < len(level); i += 2 {
			j := i + 1
			if j >= len(level) {
				j = i
			}
			next = append(next, model.HashJSON([]string{level[i], level[j]}))
		}
		pos /= 2
		level = next
	}
	return proof
}
func VerifyRoot(events []model.AuditEvent, root string) bool {
	return root != "" && RootHash(events) == root
}
