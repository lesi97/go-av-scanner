package clamscan

import (
	"bufio"
	"context"
	"errors"
	"fmt"
	"os/exec"
	"runtime"
	"strings"
	"time"

	"github.com/lesi97/go-av-scanner/internal/scanner"
	"github.com/lesi97/go-av-scanner/internal/utils"
)

const engineName = "clamdscan"

type ClamScan struct {
	logger       *utils.Logger
	BinaryPath   string
	MaxFileBytes int64
}

func New(logger *utils.Logger, binaryPath string, maxFileBytes int64) (*ClamScan, error) {
	if binaryPath == "" {
		binaryPath = engineName
	}

	if _, err := exec.LookPath(binaryPath); err != nil {
		return nil, fmt.Errorf("%s not found: %w", engineName, err)
	}

	return &ClamScan{
		logger:       logger,
		BinaryPath:   binaryPath,
		MaxFileBytes: maxFileBytes,
	}, nil
}

func (c *ClamScan) ScanFile(ctx context.Context, path string) (scanner.Result, error) {
	start := time.Now()

	// Start clamdscan and create scanners for stdout and stderr so we can read output line by line
	cmd, stdout, stderr, err := c.startCommand(ctx, path)
	if err != nil {
		c.logger.Errorf("Start Command : %v", err)
		return errorResult(start, err.Error()), err
	}

	// lines collects the stdout output so we can build a single combined string for parsing later
	// done receives the process exit result when cmd.Wait returns
	lines := make(chan string, 128)
	done := make(chan error, 1)

	// Read stdout in the background
	// We send the entire collected stdout as one chunk into lines then close the channel
	go func() {
		lines <- utils.ReadLines(stdout, c.logClamLine)
		close(lines)
	}()

	// Read stderr in the background purely for logging
	// Right now stderr is not included in output parsing, but it is useful in server logs for debugging
	go func() {
		_ = utils.ReadLines(stderr, c.logClamLine)
	}()

	// Wait for the process to finish without blocking the main goroutine
	go func() {
		done <- cmd.Wait()
	}()

	// Periodic progress log so long running scans do not look stuck
	ticker := time.NewTicker(time.Second)
	defer ticker.Stop()

	// out is a closure that blocks until stdout has been fully collected
	// This avoids races between parsing and the stdout reader goroutine
	out := collectOutput(lines)

	for {
		select {
		// If the request is cancelled, ensure the clamdscan process is killed and then return
		case <-ctx.Done():
			_ = utils.KillProcess(cmd)
			_ = <-done
			return errorResult(start, ctx.Err().Error()), ctx.Err()

		// Process finished, parse output and map exit status to our Result model
		case err := <-done:
			return c.resultFromExit(start, out(), err)

		// While running, emit a progress line occasionally
		case <-ticker.C:
			c.logger.PrintColour(true, "brightBlack", "scan running: %s elapsed\n", time.Since(start).Truncate(time.Second))
		}
	}
}

// startCommand wires up stdout and stderr pipes and starts the process
// Using exec.CommandContext links the process lifetime to the context, but we still explicitly kill on ctx.Done
func (c *ClamScan) startCommand(ctx context.Context, path string) (*exec.Cmd, *bufio.Scanner, *bufio.Scanner, error) {
	cmd := exec.CommandContext(ctx, c.BinaryPath, "--fdpass", path)

	stdoutPipe, err := cmd.StdoutPipe()
	if err != nil {
		return nil, nil, nil, err
	}

	stderrPipe, err := cmd.StderrPipe()
	if err != nil {
		return nil, nil, nil, err
	}

	if err := cmd.Start(); err != nil {
		return nil, nil, nil, err
	}

	return cmd, bufio.NewScanner(stdoutPipe), bufio.NewScanner(stderrPipe), nil
}

// logClamLine prints clamdscan output with a coloured prefix and per line colouring
// This is purely for human readability in server logs
func (c *ClamScan) logClamLine(line string) {
	prefixColour := utils.Colours["cyan"]
	reset := utils.Colours["reset"]
	lineColour := colourForClamLine(line)

	c.logger.Printf("%s%s:%s %s%s%s\n", prefixColour, engineName, reset, lineColour, line, reset)
}

// resultFromExit converts clamdscan exit conditions into a scanner.Result
// clamdscan convention is typically
// 0 means clean
// 1 means infected
// other non zero codes mean operational error
func (c *ClamScan) resultFromExit(start time.Time, output string, err error) (scanner.Result, error) {
	res := scanner.Result{
		Engine:   engineName,
		Duration: time.Since(start) / time.Millisecond,
	}

	if err == nil {
		res.Status = scanner.StatusClean
		return res, nil
	}

	var exitErr *exec.ExitError
	if errors.As(err, &exitErr) {
		code := exitErr.ExitCode()

		// Exit code 1 indicates an infected file, not an infrastructure failure
		if code == 1 {
			res.Status = scanner.StatusInfected
			res.Signature = parseSignature(output, "FOUND")
			return res, nil
		}

		// Any other exit code is treated as a scanner error
		// We extract a concise message from the tool output rather than returning noisy process errors
		res.Status = scanner.StatusError
		res.Error = fmt.Sprintf("%v", parseSignature(output, "ERROR"))
		return res, err
	}

	// If we did not get an ExitError, fall back to parsing output for a useful message
	res.Status = scanner.StatusError
	res.Error = parseSignature(output, "ERROR")
	return res, err
}

// errorResult is a small helper to ensure error responses always have engine and duration set consistently
func errorResult(start time.Time, msg string) scanner.Result {
	return scanner.Result{
		Status:   scanner.StatusError,
		Engine:   engineName,
		Duration: time.Since(start) / time.Millisecond,
		Error:    msg,
	}
}

// collectOutput gathers all text chunks sent on lines and returns a function that blocks until done
// This keeps the call site simple and prevents reading partially collected output
func collectOutput(lines <-chan string) func() string {
	var b strings.Builder
	done := make(chan struct{})

	go func() {
		for chunk := range lines {
			b.WriteString(chunk)
		}
		close(done)
	}()

	return func() string {
		<-done
		return strings.TrimSpace(b.String())
	}
}

// colourForClamLine chooses colours based on common clamdscan output patterns
// It is not part of the scanning logic, it exists to make logs readable during debugging and testing
func colourForClamLine(line string) string {
	trim := strings.TrimSpace(line)

	switch {
	case trim == "":
		return utils.Colours["brightBlack"]

	case strings.Contains(trim, "SCAN SUMMARY"):
		return utils.Colours["bold"] + utils.Colours["brightCyan"]

	case strings.HasPrefix(trim, "-----------"):
		return utils.Colours["brightBlack"]

	case strings.HasSuffix(trim, " OK"):
		return utils.Colours["green"]

	case strings.Contains(trim, " FOUND"):
		return utils.Colours["brightRed"]

	case strings.Contains(trim, " ERROR"):
		return utils.Colours["brightRed"]

	case strings.Contains(trim, "LibClamAV"):
		return utils.Colours["brightRed"]

	case strings.HasPrefix(trim, "Infected files:"):
		if strings.HasSuffix(trim, " 0") {
			return utils.Colours["green"]
		}
		return utils.Colours["brightRed"]
	}

	return utils.Colours["brightBlack"]
}

// parseSignature extracts a useful message from clamdscan output by splitting on colon
// clamdscan output often looks like
// /path/to/file: Eicar-Test-Signature FOUND
// This helper aims to return just the right side, minus the suffix marker
// On Windows, drive letters introduce extra colons, so we special case that layout
func parseSignature(output string, suffix string) string {
	parts := strings.Split(output, ":")
	if len(parts) < 2 {
		return ""
	}

	right := strings.TrimSpace(parts[len(parts)-1])

	if runtime.GOOS == "windows" && len(parts) >= 3 {
		right = strings.TrimSpace(strings.Split(parts[2], "\n")[0])
	}

	right = strings.TrimSuffix(right, suffix)
	return strings.TrimSpace(right)
}
