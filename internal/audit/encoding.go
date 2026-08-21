package audit

import (
	"encoding/json"
	"github.com/benzhi-project/a4cfdd67-9a14-48d0-ad87-ba4802711ac0/internal/model"
)

func EncodeEvent(e model.AuditEvent) []byte { b, _ := json.Marshal(e); return b }
