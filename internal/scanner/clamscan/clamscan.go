package clamscan

import (
	"context"
	"errors"
	"fmt"
	"os/exec"
	"strings"
	"time"

	"github.com/lesi97/go-av-scanner/internal/scanner"
)

type ClamScan struct {
	BinaryPath string
}

func New(binaryPath string) (ClamScan, error) {
	if binaryPath == "" {
		binaryPath = "clamscan"
	}

	if _, err := exec.LookPath(binaryPath); err != nil {
		return ClamScan{}, fmt.Errorf("clamscan not found: %w", err)
	}

	return ClamScan{BinaryPath: binaryPath}, nil
}

func (c ClamScan) ScanFile(ctx context.Context, path string) (scanner.Result, error) {
	start := time.Now()
	cmd := exec.CommandContext(ctx, c.BinaryPath, "--no-summary", path)
	outBytes, err := cmd.CombinedOutput()
	out := strings.TrimSpace(string(outBytes))

	res := scanner.Result{
		Engine:   "clamscan",
		Duration: time.Since(start) / time.Millisecond,
	}

	if err == nil {
		res.Status = scanner.StatusClean
		return res, nil
	}

	var exitErr *exec.ExitError
	if errors.As(err, &exitErr) {
		code := exitErr.ExitCode()

		if code == 1 {
			res.Status = scanner.StatusInfected
			res.Signature = parseSignature(out)
			return res, nil
		}

		res.Status = scanner.StatusError
		res.Error = out
		return res, err
	}

	res.Status = scanner.StatusError
	res.Error = out
	return res, err
}

func parseSignature(output string) string {
	parts := strings.Split(output, ":")
	if len(parts) < 2 {
		return ""
	}
	right := strings.TrimSpace(parts[len(parts)-1])
	right = strings.TrimSuffix(right, " FOUND")
	return strings.TrimSpace(right)
}
