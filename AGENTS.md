# AGENTS.md

This file provides guidance to Code Agents when working with code in this repository.

## Project Overview

go-nesgress is a Go library for hierarchical progress reporting with nested spinner displays. The name is a pun on "nested progress". It provides thread-safe progress operations with spinner animations using charmbracelet/huh.

## Build and Test Commands

```bash
# Run all tests
go test ./...

# Run tests with verbose output
go test -v ./...

# Run tests excluding integration tests (short mode)
go test -short ./...

# Run a specific test
go test -v -run TestName ./...
```

## Architecture

The library has three main components:

1. **ProgressReporter interface** (`nesgress.go`) - Defines the contract for progress reporting with methods for starting/finishing operations, persistent mode, and pause/resume support.

2. **ProgressDisplay** (`nesgress.go`) - The main implementation that manages:
   - A stack of `ProgressOperation` structs for nested hierarchy
   - Background spinner goroutines via charmbracelet/huh/spinner
   - Thread-safe state with atomic operations and mutexes
   - Terminal cursor control (hide/show)

3. **NoopProgressDisplay** (`noop.go`) - A no-op implementation for testing or disabling output.

### Thread Safety

- `synchronizedWriter` and `safeBytesBuffer` (`writer.go`) wrap io.Writer for concurrent access
- `stackMutex` protects the progress stack and active spinner
- `pauseMutex` protects pause/resume operations
- Atomic flags track operation count, cursor state, and paused state
- `spinnerWaitGroup` ensures spinner goroutines complete before state changes

### Key Patterns

- Operations use context cancellation to stop spinners
- Nested operations show hierarchical context (e.g., "Parent: Child: Grandchild")
- Completion messages show timing only for operations >100ms
- Persistent mode allows logging accomplishments that remain visible
