package model

import "time"

func Now() time.Time             { return time.Now().UTC().Truncate(time.Millisecond) }
func ValidTime(t time.Time) bool { return !t.IsZero() && t.Location() == time.UTC }
