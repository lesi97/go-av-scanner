package scanner

import (
	"context"
	"time"
)

type Status string

const (
	StatusClean   Status = "clean"
	StatusInfected Status = "infected"
	StatusError   Status = "error"
)

type Result struct {
	Status     Status        `json:"status"`
	Signature  string        `json:"signature,omitempty"`
	Engine     string        `json:"engine"`
	Duration   time.Duration `json:"duration_ms"`
	Error      string        `json:"error,omitempty"`
}

type ScanError struct {
	Result Result
}

type Scanner interface {
	ScanFile(ctx context.Context, path string) (Result, error)
}

func (e *ScanError) Error() string {
	return "file is infected"
}