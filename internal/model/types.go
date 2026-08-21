package model

import "time"

type BatchStatus string

const (
	StatusDraft      BatchStatus = "draft"
	StatusReviewing  BatchStatus = "reviewing"
	StatusCorrection BatchStatus = "correction"
	StatusFrozen     BatchStatus = "frozen"
	StatusClosed     BatchStatus = "closed"
)

type SourceState string

const (
	SourcePending         SourceState = "pending"
	SourceApproved        SourceState = "approved"
	SourceRejected        SourceState = "rejected"
	SourceNeedsCorrection SourceState = "needs_correction"
)

type Decision string

const (
	DecisionApprove    Decision = "approve"
	DecisionReject     Decision = "reject"
	DecisionCorrection Decision = "correction"
)

type DatasetBatch struct {
	BatchID       string      `json:"batch_id"`
	Title         string      `json:"title"`
	Steward       string      `json:"steward"`
	PolicyVersion string      `json:"policy_version"`
	ReleaseScope  []string    `json:"release_scope"`
	Status        BatchStatus `json:"status"`
	Version       int64       `json:"version"`
	CreatedAt     time.Time   `json:"created_at"`
	UpdatedAt     time.Time   `json:"updated_at"`
}

type SourceRecord struct {
	SourceID        string      `json:"source_id"`
	BatchID         string      `json:"batch_id"`
	Origin          string      `json:"origin"`
	Description     string      `json:"description"`
	ConsentRef      string      `json:"consent_ref"`
	Sensitivity     string      `json:"sensitivity"`
	ContentChecksum string      `json:"content_checksum"`
	EvidenceRefs    []string    `json:"evidence_refs,omitempty"`
	State           SourceState `json:"state"`
	CreatedAt       time.Time   `json:"created_at"`
}

type ReviewTask struct {
	ReviewID     string    `json:"review_id"`
	BatchID      string    `json:"batch_id"`
	SourceID     string    `json:"source_id"`
	Reviewer     string    `json:"reviewer"`
	Decision     Decision  `json:"decision"`
	Issues       []string  `json:"issues"`
	EvidenceRefs []string  `json:"evidence_refs"`
	ReviewedAt   time.Time `json:"reviewed_at"`
	Revision     int       `json:"revision"`
}

type ReleaseSnapshot struct {
	SnapshotID    string    `json:"snapshot_id"`
	BatchID       string    `json:"batch_id"`
	ManifestHash  string    `json:"manifest_hash"`
	ApprovedScope []string  `json:"approved_scope"`
	SourceCount   int       `json:"source_count"`
	FrozenAt      time.Time `json:"frozen_at"`
	FrozenBy      string    `json:"frozen_by"`
}

type AuditCertificate struct {
	CertificateID string    `json:"certificate_id"`
	BatchID       string    `json:"batch_id"`
	SnapshotID    string    `json:"snapshot_id"`
	EventHash     string    `json:"event_hash"`
	AuditEntries  int       `json:"audit_entries"`
	IssuedAt      time.Time `json:"issued_at"`
	IssuedBy      string    `json:"issued_by"`
}

type AuditEvent struct {
	Sequence     int64          `json:"sequence"`
	BatchID      string         `json:"batch_id"`
	Action       string         `json:"action"`
	Actor        string         `json:"actor"`
	Payload      map[string]any `json:"payload"`
	PreviousHash string         `json:"previous_hash"`
	EventHash    string         `json:"event_hash"`
	OccurredAt   time.Time      `json:"occurred_at"`
}

type BatchView struct {
	Batch       DatasetBatch      `json:"batch"`
	Sources     []SourceRecord    `json:"sources"`
	Reviews     []ReviewTask      `json:"reviews"`
	Snapshot    *ReleaseSnapshot  `json:"snapshot,omitempty"`
	Certificate *AuditCertificate `json:"certificate,omitempty"`
	Audit       []AuditEvent      `json:"audit"`
}
