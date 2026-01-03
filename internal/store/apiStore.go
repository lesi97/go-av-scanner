package store

import (
	"context"
	"io"

	"github.com/lesi97/go-av-scanner/internal/scanner"
	"github.com/lesi97/go-av-scanner/internal/utils"
)

type ApiStore interface {
	Health(ctx context.Context) (*string, error)
	Scan(ctx context.Context, file io.Reader) (*scanner.Result, error)
	MaxUploadBytes() int64
}

type DbApiStore struct {
	logger *utils.Logger
	scanner scanner.Scanner
	sem     chan struct{}
	maxUploadBytes int64
}

func NewApiStore(logger *utils.Logger, sc scanner.Scanner, maxUploadBytes int64,) *DbApiStore {
	return &DbApiStore{
		logger: logger,
		scanner: sc,
		sem: make(chan struct{}, 20),
		maxUploadBytes: maxUploadBytes,
	}
}

func (s *DbApiStore) MaxUploadBytes() int64 {
	return s.maxUploadBytes
}