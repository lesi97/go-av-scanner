package api_test

import (
	"bytes"
	"context"
	"errors"
	"io"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/lesi97/go-av-scanner/internal/api"
	"github.com/lesi97/go-av-scanner/internal/scanner"
	"github.com/lesi97/go-av-scanner/internal/utils"
)

type fakeStore struct {
	res *scanner.Result
	err error
}

func (f fakeStore) Scan(ctx context.Context, r io.Reader) (*scanner.Result, error) {
	return f.res, f.err
}

func (f fakeStore) Health(ctx context.Context) (*string, error) {
	return nil, nil
}


func TestHandleScan_ContentField(t *testing.T) {
	logger := utils.NewColourLogger("brightMagenta")
	h := api.NewApiHandler(logger, fakeStore{
		res: &scanner.Result{Status: scanner.StatusClean, Engine: "fake"},
	})

	req := httptest.NewRequest(http.MethodPost, "/scan", strings.NewReader(""))
	req.Header.Set("Content-Type", "multipart/form-data; boundary=BOUNDARY")

	body := "--BOUNDARY\r\n" +
		"Content-Disposition: form-data; name=\"content\"\r\n\r\n" +
		"hello\r\n" +
		"--BOUNDARY--\r\n"
	req.Body = io.NopCloser(strings.NewReader(body))

	w := httptest.NewRecorder()
	h.HandleScan(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %v body=%v", w.Code, w.Body.String())
	}
}

func TestHandleScan_InfectedMapsTo422(t *testing.T) {
	logger := utils.NewColourLogger("brightMagenta")
	h := api.NewApiHandler(logger, fakeStore{
		res: &scanner.Result{Status: scanner.StatusInfected, Signature: "Eicar-Test-Signature", Engine: "fake"},
		err: &scanner.ScanError{Result: scanner.Result{Status: scanner.StatusInfected, Signature: "Eicar-Test-Signature", Engine: "fake"}},
	})

	body := &bytes.Buffer{}
	mw := multipart.NewWriter(body)
	_ = mw.WriteField("content", "eicar")
	_ = mw.Close()

	req := httptest.NewRequest(http.MethodPost, "/scan", body)
	req.Header.Set("Content-Type", mw.FormDataContentType())

	w := httptest.NewRecorder()
	h.HandleScan(w, req)

	if w.Code != http.StatusUnprocessableEntity {
		t.Fatalf("expected 422, got %v body=%v", w.Code, w.Body.String())
	}
}

func TestHandleScan_InternalErrorMapsTo500(t *testing.T) {
	logger := utils.NewColourLogger("brightMagenta")
	h := api.NewApiHandler(logger, fakeStore{
		res: nil,
		err: errors.New("boom"),
	})

	body := &bytes.Buffer{}
	mw := multipart.NewWriter(body)
	_ = mw.WriteField("content", "hello")
	_ = mw.Close()

	req := httptest.NewRequest(http.MethodPost, "/scan", body)
	req.Header.Set("Content-Type", mw.FormDataContentType())

	w := httptest.NewRecorder()
	h.HandleScan(w, req)

	if w.Code != http.StatusInternalServerError {
		t.Fatalf("expected 500, got %v body=%v", w.Code, w.Body.String())
	}
}
