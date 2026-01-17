---
name: go-nesgress-dev
description: Development guide for the go-nesgress library. Use when working on the hierarchical progress reporting library, implementing features, testing, or maintaining the codebase. Covers testing strategies with noop pattern, development workflows, design principles, and maintenance practices. For general Go coding conventions and testing patterns, use the go-dev skill.
---

# go-nesgress Development

Development guide for the go-nesgress Go library - a hierarchical progress reporting system with nested spinner displays.

## Project Overview

go-nesgress is a library that provides npm-style progress reporting with:
- Hierarchical parent-child progress operations
- Thread-safe concurrent operation support
- Animated spinners using charmbracelet/huh
- Success/failure indicators with timing
- Persistent mode showing accomplishments
- Pause/resume for interactive prompts

## Project Structure

```
go-nesgress/
├── nesgress.go           # Core ProgressDisplay implementation
├── noop.go              # NoopProgressDisplay implementation
├── writer.go            # Thread-safe writer utilities
├── doc.go               # Package documentation
├── nesgress_test.go     # Test suite
├── go.mod               # Dependencies (charmbracelet/huh, lipgloss)
├── Taskfile.yml         # Development tasks
└── .golangci.yml        # Linter configuration
```

## Development Commands

Use Taskfile for common operations:

| Command | Purpose |
|---------|---------|
| `task test` | Run tests with race detection |
| `task fmt` | Format code with go fmt |
| `task lint` | Run golangci-lint and typos |
| `task check` | Run tests + lint together |
| `task cov` | Generate coverage report |
| `task bench` | Run benchmarks |
| `task doc` | Render pkg docs locally with pkgsite |
| `task tidy` | Tidy dependencies |
| `task clean` | Remove build artifacts |

## Testing Strategies

### Using NoopProgressDisplay for Testing

**Key insight:** Tests should use the noop implementation to avoid terminal output and spinner animations. This is a critical pattern for this library.

The noop implementation provides a drop-in replacement that:
- Succeeds silently for all operations
- Returns no errors
- Produces no output
- Requires no cleanup or synchronization

Use `NewNoopProgressDisplay()` in test code and anywhere progress display should be disabled (CI environments, silent mode, etc.). This avoids the complexity of mocking while maintaining type safety.

### Testing the Actual Display Implementation

When you need to test the actual display behavior (formatting, output, timing):
- Use a bytes.Buffer as the output writer
- Capture and verify the output string
- Remember to call Close() to clean up goroutines
- Tests with actual display are slower due to spinner animations

### Concurrency Testing

The library is designed for concurrent use. When testing or adding features:
- Run tests with `-race` flag (enabled by default in Taskfile)
- Test multiple goroutines calling progress methods simultaneously
- Verify no data races or deadlocks
- Consider goroutine lifecycle and cleanup

## Design Principles

### Two Implementation Pattern

The library provides exactly two implementations - full and noop. This is intentional:
- **Full implementation** - Complete functionality with terminal output
- **Noop implementation** - Silent, no-op version for testing

This pattern enables dependency injection without mocking. Consumers can inject the appropriate implementation based on context (production vs testing, verbose vs silent).

### Thread Safety by Default

All public methods are thread-safe. Users should never need external synchronization when calling library methods, even from multiple goroutines.

When adding features, ensure thread safety through:
- Synchronized wrappers for I/O operations
- Appropriate locks for shared state
- Proper goroutine lifecycle management

### Zero Configuration Philosophy

The library works out of the box with sensible defaults:
- No configuration files or setup required
- No initialization beyond calling the constructor
- Minimal API surface area
- Reasonable defaults for all behavior (timing display, spinner speed, etc.)

Avoid adding configuration unless absolutely necessary. Prefer sensible defaults.

## Adding Features

### Adding New Progress Methods

When adding new methods to the library:
1. Add to the main interface in nesgress.go
2. Implement the full version in ProgressDisplay
3. Implement the silent version in NoopProgressDisplay (usually just `return nil`)
4. Add comprehensive tests for the new functionality
5. Update doc.go with examples
6. Update README.md if user-facing

Both implementations must always stay in sync.

### Modifying Display Behavior

Display modifications typically involve:
- Message formatting and context display
- Success/failure indicator symbols (✓, ✗)
- Timing display thresholds and formatting
- Styling via lipgloss (colors, emphasis)

Test display changes manually with terminal output, not just unit tests. Spinner behavior and visual appearance matter.

### Thread Safety Considerations

When modifying the library:
- Use existing synchronized types for I/O operations
- Protect all shared state with appropriate locking mechanisms
- Always test with `-race` flag enabled
- Consider goroutine cleanup and cancellation paths
- Avoid introducing global state

## Common Development Patterns

### Dependency Injection

The library is designed to be injected as an interface dependency. Consumers should accept the interface type and inject either the full or noop implementation based on runtime needs.

### Error Handling Flow

The standard flow is: Start an operation, perform work, then either Fail (with error context) or Finish (with success message). The library handles display state transitions automatically.

### Persistent Mode Usage

Persistent mode is for long-running operations with multiple visible accomplishments. Each accomplishment gets logged and remains visible while the operation continues. Use for deploy pipelines, multi-step installations, or batch processing where users benefit from seeing progressive completion.

## Dependencies

- **charmbracelet/huh/spinner** - Spinner animations
- **charmbracelet/lipgloss** - Terminal styling

Keep dependencies minimal. This is a library, not an application. New dependencies require strong justification.

## Maintenance Notes

### Linting

The `.golangci.yml` is comprehensive with most linters enabled. Check comments in the file for disabled linters and their rationale. Don't modify linter settings without understanding why they're configured as they are.

### Documentation Synchronization

Keep these in sync:
- **doc.go** - Package documentation with examples
- **README.md** - User-facing documentation
- **SKILL.md** - Development guidance (this file)

When adding features or changing behavior, all three need updates.

### Releases

Project uses Go modules. Follow semantic versioning:
- Patch: Bug fixes, performance improvements
- Minor: New features, backwards compatible
- Major: Breaking API changes

### Testing Before Release

Before releasing:
1. Run `task check` (tests + linting)
2. Run `task bench` and verify no performance regressions
3. Manually test visual display in terminal
4. Verify README examples still work
5. Update CHANGELOG if maintained
