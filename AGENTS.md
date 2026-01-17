# AGENTS.md

This file provides context for Code Agents working with this repository.

## What Problem Does This Solve?

Terminal-based applications need to communicate progress to users during long-running operations. Simple approaches have limitations:

**Print-based progress** (e.g., "Downloading...", "Downloaded") creates cluttered output when operations are nested. For example, a deployment might involve downloading, compiling, and uploading - each with their own sub-steps. Print-based progress loses the hierarchical context.

**Spinner libraries** improve visual appeal but typically handle only flat progress. They don't naturally support nested operations where you want to show "Deploying: Building: Compiling main.go" as context.

**Custom solutions** end up reimplementing the same patterns: managing a stack of operations, handling terminal control codes, coordinating goroutines for animations, and ensuring thread safety. This is error-prone and distracts from actual application logic.

go-nesgress solves this by providing hierarchical progress reporting out of the box. It maintains context across nested operations, handles all terminal complexity, and works safely from concurrent goroutines.

## Design Goals

**Hierarchical by default** - Progress operations naturally nest. The library's core model is a stack that tracks parent-child relationships and displays them as context.

**Zero configuration** - No setup, no initialization, no decisions. Call the constructor and it works with sensible defaults.

**Thread-safe** - All methods work safely from any goroutine without requiring external synchronization. Concurrent progress reporting just works.

**Non-invasive** - Progress reporting shouldn't affect application logic. If display fails, the application continues normally.

**Clean visuals** - Spinners animate smoothly, completions show timing, hierarchy provides context. The goal is npm-style polish in Go applications.

## When to Use This Library

**Multi-step operations** where users benefit from seeing progress context (installers, deployment tools, build systems).

**Nested workflows** where inner operations need to report progress while maintaining outer context (package managers, migration tools).

**Concurrent operations** where multiple goroutines report progress simultaneously (parallel downloads, batch processing).

**Long-running CLI tools** where visual feedback prevents users from thinking the program has hung.

## When NOT to Use This Library

**Silent background services** - If there's no user watching, progress display adds overhead without benefit.

**High-frequency operations** - Sub-100ms operations don't benefit from progress display. The library filters these out, but if most operations are trivial, the library adds unnecessary complexity.

**Non-terminal output** - The library assumes terminal output. It won't work well if output is redirected to files or parsed by other programs (though the noop implementation can handle this case).

**GUIs or web interfaces** - This is specifically for terminal/CLI applications.

## Technical Resources

- **Architecture and patterns:** [docs/architecture.md](docs/architecture.md)
- **Build commands:** Use `task` (see Taskfile.yml) - `task test`, `task lint`, `task check`

## Philosophy

This is a **library, not a framework**. It does one thing (hierarchical progress reporting) and does it well. It has minimal dependencies, a small API surface, and stays out of the way.

The library should feel natural to Go developers. It uses standard patterns (interfaces, context cancellation, io.Writer) and follows Go conventions. There's no magic, no global state, no complex configuration.

Progress reporting is supporting functionality, not the main show. The library is designed so applications can integrate it easily and forget about it. It works quietly in the background, providing value without demanding attention.
