package store_test

import (
	"bytes"
	"context"
	"errors"
	"io"
	"testing"

	"github.com/lesi97/go-av-scanner/internal/scanner"
	"github.com/lesi97/go-av-scanner/internal/store"
	"github.com/lesi97/go-av-scanner/internal/utils"
)

type fakeScanner struct {
	result scanner.Result
	err    error
}

func (f fakeScanner) ScanFile(ctx context.Context, path string) (scanner.Result, error) {
	return f.result, f.err
}

const MaxUploadBytes int64 = 64 << 20

func TestScan_ReturnsErrorOnNilReader(t *testing.T) {
	logger := utils.NewColourLogger("brightMagenta")
	s := store.NewApiStore(logger, fakeScanner{}, MaxUploadBytes)

	_, err := s.Scan(context.Background(), nil)
	if err == nil {
		t.Fatalf("expected error, got nil")
	}
}

func TestScan_ReturnsCleanResult(t *testing.T) {
	fs := fakeScanner{
		result: scanner.Result{Status: scanner.StatusClean, Engine: "fake"},
	}
	logger := utils.NewColourLogger("brightMagenta")
	s := store.NewApiStore(logger, fs, MaxUploadBytes)

	res, err := s.Scan(context.Background(), bytes.NewReader([]byte("hello")))
	if err != nil {
		t.Fatalf("expected nil error, got %v", err)
	}
	if res == nil {
		t.Fatalf("expected result, got nil")
	}
	if res.Status != scanner.StatusClean {
		t.Fatalf("expected clean, got %v", res.Status)
	}
}

func TestScan_ReturnsInfectedResult(t *testing.T) {
	fs := fakeScanner{
		result: scanner.Result{Status: scanner.StatusInfected, Signature: "Eicar-Test-Signature", Engine: "fake"},
	}
	logger := utils.NewColourLogger("brightMagenta")
	s := store.NewApiStore(logger, fs, MaxUploadBytes)

	res, err := s.Scan(context.Background(), bytes.NewReader([]byte("eicar")))
	if err != nil {
		t.Fatalf("expected nil error, got %v", err)
	}
	if res == nil {
		t.Fatalf("expected result, got nil")
	}
	if res.Status != scanner.StatusInfected {
		t.Fatalf("expected infected, got %v", res.Status)
	}
	if res.Signature != "Eicar-Test-Signature" {
		t.Fatalf("expected signature, got %v", res.Signature)
	}
}


func TestScan_PropagatesScannerError(t *testing.T) {
	fs := fakeScanner{
		result: scanner.Result{Status: scanner.StatusError, Engine: "fake"},
		err:    io.EOF,
	}
	logger := utils.NewColourLogger("brightMagenta")
	s := store.NewApiStore(logger, fs, MaxUploadBytes)

	res, err := s.Scan(context.Background(), bytes.NewReader([]byte("hello")))
	if res == nil {
		t.Fatalf("expected result, got nil")
	}
	if !errors.Is(err, io.EOF) {
		t.Fatalf("expected EOF, got %v", err)
	}
}

func TestMaxUploadBytes_ReturnsConfiguredValue(t *testing.T) {
	logger := utils.NewColourLogger("brightMagenta")
	s := store.NewApiStore(logger, fakeScanner{}, MaxUploadBytes)

	if s.MaxUploadBytes() != MaxUploadBytes {
		t.Fatalf("expected %d, got %d", MaxUploadBytes, s.MaxUploadBytes())
	}
}
