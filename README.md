# 🥷 Bandito

This README provides detailed instructions on how to build the Bandito application for different architectures, specifically Darwin ARM (Apple Silicon) and AMD64.

## Prerequisites

Before you begin, ensure you have the following installed:

- Go 1.16 or higher


## Building

Build the app and outputs the binary to the ./bin/ directory.

### Darwin arm64 (Apple Silicon)

```bash
GOOS=darwin GOARCH=arm64 go build -o ./bin/bandito-darwin-arm64
```

### Darwin amd64 (Intel)

```bash
GOOS=darwin GOARCH=amd64 go build -o ./bin/bandito-darwin-amd64
```