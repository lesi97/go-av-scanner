# go-av-scanner

A lightweight Go HTTP service for scanning uploaded content with ClamAV

Designed to be simple, testable, and easy to run locally or in Docker

## Features

-   `POST /scan` endpoint
-   Accepts file uploads or raw content
-   ClamAV scanning via `clamscan`
-   Single container Docker setup
-   Integration tests that auto skip when ClamAV is not available

## API

### POST /scan

Send either:

-   `file` as a multipart file upload
-   `content` as a multipart text field

Example:

```bash
curl -X POST http://localhost:8080/scan \ -F "file=@./example.txt"
```

## Docker

### Build

```bash
docker build -t go-av-scanner .
```

### Run

```bash
docker run --rm -p 8080:8080 -v clamdb:/var/lib/clamav go-av-scanner
```

Notes:

-   First run may take longer while ClamAV signatures download
-   The `clamdb` volume persists the virus database between restarts

## Local development

Run the service:

```bash
go run ./cmd/scanner-api
```

ClamAV must be installed locally if not using Docker

## Tests

Run all tests:

```bash
go test ./...
```

Integration tests:

-   Automatically skip if `clamscan` is not available
-   Use `CLAMSCAN_PATH` to point to a specific binary if needed

Example:

```bash
CLAMSCAN_PATH=/usr/bin/clamscan go test ./... -count=1
```

## Installing ClamAV

### Windows

-   Install ClamAV for Windows
-   Ensure `clamscan.exe` exists, commonly at:
    ```
    C:\Program Files\ClamAV\clamscan.exe
    ```
-   Add `C:\Program Files\ClamAV` to PATH
-   Restart your terminal
-   Verify:

```powershell
where.exe clamscan clamscan --version
```

### Linux

```bash
sudo apt-get install -y clamav sudo freshclam
```

### macOS

```bash
brew install clamav freshclam
```

## Support

If you find this project useful, you can support its development via GitHub Sponsors

## Licence

MIT
