package clamscan_test

import (
	"context"
	"os"
	"os/exec"
	"strings"
	"testing"
	"time"

	"github.com/lesi97/go-av-scanner/internal/scanner"
	"github.com/lesi97/go-av-scanner/internal/scanner/clamscan"
)

func TestClamScan_EicarDetected(t *testing.T) {
	bin := os.Getenv("CLAMSCAN_PATH")
	if bin == "" {
		bin = "clamscan"
	}

	if _, err := exec.LookPath(bin); err != nil {
		t.Skip("clamscan not found, skipping integration test")
	}

	sc, _ := clamscan.New(bin)

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	tmp, err := os.CreateTemp("", "eicar-*")
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
