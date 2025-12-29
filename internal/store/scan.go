package store

import (
	"context"
	"fmt"
	"io"
	"os"

	"github.com/lesi97/go-av-scanner/internal/scanner"
)

func (s *DbApiStore) Scan(ctx context.Context, r io.Reader) (*scanner.Result, error) {
	if r == nil {
		return nil, fmt.Errorf("missing request body")
	}

	tmp, err := os.CreateTemp("", "upload-*")
	if err != nil {
		return nil, fmt.Errorf("failed to create temp file: %w", err)
	}

	tmpPath := tmp.Name()
	defer func() { _ = os.Remove(tmpPath) }()

	_, err = io.Copy(tmp, r)
	if err != nil {
		_ = tmp.Close()
		return nil, fmt.Errorf("failed to read upload: %w", err)
	}

	err = tmp.Close()
	if err != nil {
		return nil, fmt.Errorf("failed to close temp file: %w", err)
	}
	
	result, scanErr := s.scanner.ScanFile(ctx, tmpPath)
	if scanErr != nil {
		return &result, scanErr
	}

	if result.Status == scanner.StatusInfected {
		return &result, &scanner.ScanError{Result: result}
	}

	return &result, nil
}
