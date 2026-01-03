package store

import (
	"context"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"time"

	"github.com/lesi97/go-av-scanner/internal/scanner"
)

func (s *DbApiStore) Scan(ctx context.Context, r io.Reader) (*scanner.Result, error) {

	select {
	case s.sem <- struct{}{}:
		defer func() { <-s.sem }()
	case <-ctx.Done():
		return nil, ctx.Err()
	}

	if r == nil {
		return nil, fmt.Errorf("missing request body")
	}

	tmpDir, err := getScanTmpDir()
	if err != nil {
		return nil, err
	} 

	tmp, err := os.CreateTemp(tmpDir, "upload-*")
	if err != nil {
		return nil, fmt.Errorf("failed to create temp file: %w", err)
	}

	tmpPath := tmp.Name()
	defer func() { _ = os.Remove(tmpPath) }()

	t0 := time.Now()
	_, err = io.Copy(tmp, r)
	s.logger.LogExecutionTime("copyToTmp", t0)
	if err != nil {
		_ = tmp.Close()
		return nil, fmt.Errorf("failed to read upload: %w", err)
	}

	err = tmp.Close()
	if err != nil {
		return nil, fmt.Errorf("failed to close temp file: %w", err)
	}
	
	t0 = time.Now()
	result, scanErr := s.scanner.ScanFile(ctx, tmpPath)
	s.logger.LogExecutionTime("clamdscan", t0)
	if scanErr != nil {
		return &result, scanErr
	}

	return &result, nil
}

func getScanTmpDir() (string, error) {
	if v := strings.TrimSpace(os.Getenv("AV_SCAN_TMPDIR")); v != "" {
		if err := os.MkdirAll(v, 0o700); err != nil {
			return "", fmt.Errorf("failed to create AV_SCAN_TMPDIR: %w", err)
		}
		return v, nil
	}

	if runtime.GOOS == "windows" {
		base := `C:\av-test-tmp`
		if err := os.MkdirAll(base, 0o700); err != nil {
			return "", fmt.Errorf("failed to create windows tmp dir: %w", err)
		}
		return base, nil
	}

	base := filepath.Join(os.TempDir(), "av-scan")
	if err := os.MkdirAll(base, 0o700); err != nil {
		return "", fmt.Errorf("failed to create tmp dir: %w", err)
	}
	return base, nil
}