package model

type Transition struct {
	From   BatchStatus
	Action string
	To     BatchStatus
}

var transitions = []Transition{{StatusDraft, "source_added", StatusReviewing}, {StatusReviewing, "source_reviewed", StatusReviewing}, {StatusReviewing, "correction_required", StatusCorrection}, {StatusCorrection, "correction_resubmitted", StatusReviewing}, {StatusReviewing, "batch_frozen", StatusFrozen}, {StatusFrozen, "batch_closed", StatusClosed}}

func AllowedTransition(from BatchStatus, action string) (BatchStatus, bool) {
	for _, t := range transitions {
		if t.From == from && t.Action == action {
			return t.To, true
		}
	}
	return "", false
}
func TransitionList() []Transition       { return append([]Transition(nil), transitions...) }
func IsTerminal(status BatchStatus) bool { return status == StatusClosed }
func IsWritable(status BatchStatus) bool { return status != StatusClosed }
func PendingSources(sources []SourceRecord) []string {
	out := []string{}
	for _, s := range sources {
		if s.State != SourceApproved {
			out = append(out, s.SourceID)
		}
	}
	return out
}
func CorrectionSources(sources []SourceRecord) []string {
	out := []string{}
	for _, s := range sources {
		if s.State == SourceNeedsCorrection || s.State == SourceRejected {
			out = append(out, s.SourceID)
		}
	}
	return out
}
func ScopeAllows(scope []string, value string) bool {
	for _, x := range scope {
		if x == value {
			return true
		}
	}
	return false
}
