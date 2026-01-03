package clamscan_test

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"testing"
	"time"

	"github.com/lesi97/go-av-scanner/internal/scanner"
	"github.com/lesi97/go-av-scanner/internal/scanner/clamscan"
	"github.com/lesi97/go-av-scanner/internal/utils"
)

func TestClamScan_EicarDetected(t *testing.T) {
	logger := utils.NewColourLogger("brightMagenta")
	bin := os.Getenv("CLAMSCAN_PATH")
	if bin == "" {
		bin = "clamdscan"
	}

	if _, err := exec.LookPath(bin); err != nil {
		t.Skip("clamdscan not found, skipping integration test")
	}

	sc, _ := clamscan.New(logger, bin, 64 << 20)

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	tmpDir, err := getScanTmpDir()
	if err != nil {
		t.Fatalf("temp: %v", err)
	} 

	tmp, err := os.CreateTemp(tmpDir, "eicar-*")
	if err != nil {
		t.Fatalf("temp: %v", err)
	}
	path := tmp.Name()
	defer func() { _ = os.Remove(path) }()

	_, _ = tmp.WriteString("X5O!P%@AP[4\\PZX54(P^)7CC)7}$EICAR-STANDARD-ANTIVIRUS-TEST-FILE!$H+H*")
	_ = tmp.Close()

	res, err := sc.ScanFile(ctx, path)
	if err != nil {
		t.Fatalf("scan err: %v result=%+v", err, res)
	}
	if res.Status != scanner.StatusInfected {
		t.Fatalf("expected infected got %v", res.Status)
	}
	if !strings.Contains(res.Signature, "Eicar") {
		t.Fatalf("expected eicar signature got %v", res.Signature)
	}
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
