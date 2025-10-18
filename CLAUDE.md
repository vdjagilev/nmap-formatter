# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

nmap-formatter is a CLI tool that converts NMAP XML scan output into multiple formats: HTML, CSV, JSON, Markdown, Excel, SQLite, Graphviz (dot), and D2 lang. It's built in Go and uses a workflow-based architecture with pluggable formatters.

## Development Commands

### Build and Test
```bash
# Build the project
go build -v ./

# Run all tests
go test -v ./...

# Run tests with coverage
go test ./... -race -coverprofile=coverage.txt -covermode=atomic

# Run linter (same as CI)
golangci-lint run --timeout 10m
```

### Running the Application
```bash
# Basic usage
./nmap-formatter [html|csv|md|json|dot|sqlite|excel|d2] [path-to-nmap.xml]

# From stdin
cat example.xml | ./nmap-formatter json

# With output file
./nmap-formatter html scan.xml -f output.html
```

### Testing Individual Packages
```bash
# Test only the formatter package
go test -v ./formatter/...

# Test only the cmd package
go test -v ./cmd/...

# Run a specific test
go test -v -run TestWorkflow_Execute ./formatter/
```

## Architecture

### Core Components

**Entry Point (main.go)**: Minimal entry point that delegates to `cmd.Execute()`.

**Command Layer (cmd/)**: Uses Cobra for CLI argument parsing and flag management. The `root.go` file defines all CLI flags and initializes the `MainWorkflow`.

**Formatter Package (formatter/)**: Contains the core business logic:

- **Workflow Pattern**: `MainWorkflow` (implements `Workflow` interface) orchestrates the conversion pipeline:
  1. Parse NMAP XML input into `NMAPRun` struct
  2. Apply filter expressions using the expr-lang library
  3. Build `TemplateData` with scan results and output options
  4. Delegate to format-specific `Formatter` implementation
  5. Write output to file or STDOUT

- **Data Models**: NMAP XML is mapped to Go structs:
  - `NMAPRun`: Root scan metadata (scanner version, args, timestamps)
  - `Host`: Individual scanned host (addresses, hostnames, status, OS, trace, ports)
  - `Port`: Port details (ID, protocol, state, service, scripts)
  - These structs use XML struct tags for unmarshaling

- **Formatter Interface**: All output formats implement `Formatter` interface:
  ```go
  type Formatter interface {
      Format(td *TemplateData, templateContent string) error
      defaultTemplateContent() string
  }
  ```
  The factory function `formatter.New(config)` returns the appropriate formatter based on config.

- **Format Implementations**:
  - `HTMLFormatter` / `MarkdownFormatter`: Use Go templates from `resources/templates/`
  - `JSONFormatter`: Uses `encoding/json`
  - `CSVFormatter`: Custom CSV generation
  - `ExcelFormatter`: Uses excelize library
  - `SqliteFormatter`: Writes to SQLite using multiple repository types (hosts, ports, OS, scans)
  - `DotFormatter`: Generates Graphviz syntax
  - `D2LangFormatter`: Generates D2 diagram language

- **Filtering**: The `expr.go` file integrates the expr-lang library to filter hosts based on expressions. The `--filter` flag and `--skip-down-hosts` both use this mechanism.

### Configuration Flow

1. `cmd/root.go` initializes `formatter.Config` with all CLI flags
2. Config contains `OutputOptions` with format-specific settings (e.g., `HTMLOptions`, `MarkdownOptions`)
3. Config is passed to `MainWorkflow.SetConfig()`, then to formatter instances
4. Each formatter reads its relevant options from `config.OutputOptions`

### Important Patterns

- **Input/Output Handling**: Supports both file paths and STDIN/STDOUT. The workflow sets up `config.InputFileConfig.Source` and `config.Writer` before execution.

- **Template Customization**: HTML and Markdown formatters support custom templates via `--html-use-template` and `--md-use-template` flags. Templates are loaded by `TemplateContent()` function.

- **SQLite Architecture**: Uses repository pattern with separate repositories for hosts, ports, OS, and scans. The `sqlite_db.go` manages schema creation.

- **Version Management**: Version is hardcoded in `cmd/root.go` as `const VERSION`. This should be updated for releases.

## Testing

All major components have corresponding `_test.go` files. Tests typically:
- Use sample XML fixtures for parsing tests
- Mock `io.Writer` for output validation
- Test filter expressions with various host/port conditions
- Validate format-specific output structure

The CI runs tests on Go 1.24.x across Linux, macOS, and Windows.

## Module Information

- Module path: `github.com/vdjagilev/nmap-formatter/v3`
- Go version: 1.24
- Uses go modules for dependency management

## Key Dependencies

- `github.com/spf13/cobra`: CLI framework
- `github.com/expr-lang/expr`: Expression language for filtering
- `github.com/xuri/excelize/v2`: Excel file generation
- `github.com/mattn/go-sqlite3`: SQLite driver (CGO-based)
- `oss.terrastruct.com/d2`: D2 diagram language library
- `golang.org/x/net`: Used for HTML parsing in certain formatters
