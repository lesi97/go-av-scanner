package store

import (
	"context"
	"io"
	"log"

	"github.com/lesi97/go-av-scanner/internal/scanner"
)

type ApiStore interface {
	Health(ctx context.Context) (*string, error)
	Scan(ctx context.Context, file io.Reader) (*scanner.Result, error)
}

type DbApiStore struct {
	logger *log.Logger
	scanner scanner.Scanner
	sem     chan struct{}
}

func NewApiStore(logger *log.Logger, sc scanner.Scanner) *DbApiStore {
	return &DbApiStore{
		logger: logger,
		scanner: sc,
		sem: make(chan struct{}, 20),
	}
}
