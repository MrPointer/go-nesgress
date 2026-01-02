# go-nesgress

A Go library for hierarchical progress reporting with nested spinner displays.

The name "nesgress" is a pun on "nested progress", reflecting the library's ability to display hierarchical, nested progress operations in the terminal.

## Features

- Hierarchical progress display with parent-child relationships
- Thread-safe concurrent operation support
- Spinner animations using [charmbracelet/huh](https://github.com/charmbracelet/huh)
- Success/failure indicators with timing information
- Persistent mode for long-running operations with accomplishments
- Pause/resume support for interactive prompts

## Installation

```bash
go get github.com/MrPointer/go-nesgress
```

## Usage

### Basic Progress

```go
package main

import (
    "os"
    "time"

    "github.com/MrPointer/go-nesgress"
)

func main() {
    display := nesgress.NewProgressDisplay(os.Stdout)
    defer display.Close()

    display.Start("Installing packages")
    time.Sleep(2 * time.Second) // Simulate work
    display.Finish("Packages installed")
}
```

Output:
```
✓ Packages installed (took 2s)
```

### Nested Operations

```go
display.Start("Setting up environment")

    display.Start("Downloading dependencies")
    time.Sleep(1 * time.Second)
    display.Finish("Dependencies downloaded")

    display.Start("Compiling")
    time.Sleep(500 * time.Millisecond)
    display.Finish("Compiled")

display.Finish("Environment ready")
```

While running, shows hierarchical context:
```
⠋ Setting up environment: Downloading dependencies
```

### Handling Errors

```go
display.Start("Connecting to database")
err := connectToDatabase()
if err != nil {
    display.Fail("Connection failed", err)
    return
}
display.Finish("Connected")
```

Output on failure:
```
✗ Connection failed
  Error: connection refused
```

### Persistent Mode

For long-running operations where you want to show intermediate accomplishments:

```go
display.StartPersistent("Deploying application")

display.LogAccomplishment("Built container")
display.LogAccomplishment("Pushed to registry")
display.LogAccomplishment("Updated load balancer")

display.FinishPersistent("Deployment complete")
```

Output:
```
   ✓ Built container
   ✓ Pushed to registry
   ✓ Updated load balancer
✓ Deploying application (took 45s)
```

### Pause/Resume for Interactive Input

When you need to prompt for user input:

```go
display.Start("Processing files")

// Need to ask user something
display.Pause()
fmt.Print("Continue? [y/n]: ")
// ... get input ...
display.Resume()

display.Finish("Files processed")
```

### Noop Implementation

For testing or when progress display should be disabled:

```go
var display nesgress.ProgressReporter

if verbose {
    display = nesgress.NewProgressDisplay(os.Stdout)
} else {
    display = nesgress.NewNoopProgressDisplay()
}

// Same API, no output when using Noop
display.Start("Working...")
display.Finish("Done")
```

## API Reference

### ProgressReporter Interface

```go
type ProgressReporter interface {
    io.Closer
    Start(message string) error
    Update(message string) error
    Finish(message string) error
    Fail(message string, err error) error
    StartPersistent(message string) error
    LogAccomplishment(message string) error
    FinishPersistent(message string) error
    FailPersistent(message string, err error) error
    Clear() error
    Pause() error
    Resume() error
    IsActive() bool
    IsPaused() bool
}
```

### Functions

- `NewProgressDisplay(output io.Writer) *ProgressDisplay` - Create a new progress display
- `NewNoopProgressDisplay() *NoopProgressDisplay` - Create a no-op progress display

## Dependencies

- [github.com/charmbracelet/huh](https://github.com/charmbracelet/huh) - Spinner animations
- [github.com/charmbracelet/lipgloss](https://github.com/charmbracelet/lipgloss) - Terminal styling

## License

MIT License - see [LICENSE](LICENSE) for details.
