# go-nesgress Architecture

This document describes the architectural patterns and design decisions in go-nesgress, a hierarchical progress reporting library.

## Core Model: Hierarchical Progress

Progress operations are organized in a parent-child hierarchy using a stack data structure:

```
Operation 1 (root)
├── Operation 2 (child)
│   ├── Operation 3 (grandchild)
│   └── Operation 4 (grandchild)
└── Operation 5 (child)
```

**Key behaviors:**
- Only the deepest (most recent) operation shows an active spinner
- Spinner messages display the full hierarchy path: "Parent: Child: Grandchild"
- Completing a child operation returns focus to its parent
- Each operation tracks its start time for duration display

**Why hierarchical:**
- Provides context for nested operations (e.g., "Installing: Downloading dependencies")
- Matches natural structure of complex operations
- Allows granular progress reporting without losing big picture

## Thread Safety Strategy

All public methods are thread-safe by default. Multiple goroutines can safely call progress methods concurrently without external synchronization.

**Thread safety mechanisms:**

1. **Synchronized Writers** - Wrap io.Writer with mutex protection for concurrent writes
2. **Stack Protection** - Mutex guards the progress stack during push/pop operations
3. **Atomic Flags** - Used for simple state (operation count, pause state, cursor visibility)
4. **Pause Protection** - Separate mutex for pause/resume to avoid deadlocks
5. **WaitGroup Coordination** - Ensures spinner goroutines complete before state changes

**Why thread-safe by default:**
- Library users shouldn't think about synchronization
- Enables natural concurrent usage patterns
- Prevents subtle bugs in multi-threaded applications

## Spinner Lifecycle Pattern

Each active operation spawns a background goroutine to run the spinner animation:

1. **Start** - Create context with cancellation, spawn goroutine, run spinner
2. **Update** - Cancel previous spinner, spawn new one with updated message
3. **Finish/Fail** - Cancel spinner, print completion message, clean up
4. **Pause** - Cancel all spinners, clear terminal line
5. **Resume** - Restart spinner for active operation

**Context cancellation** is the primary mechanism for stopping spinners. When an operation completes or pauses, its context is cancelled, causing the spinner goroutine to exit cleanly.

**WaitGroup synchronization** ensures spinner goroutines fully terminate before modifying state. This prevents race conditions where a spinner might write to output after its operation has completed.

**Why goroutines for spinners:**
- Non-blocking progress operations
- Smooth animations without polling in application code
- Clean cancellation via context
- Natural fit for concurrent Go code

## Persistent Mode Pattern

Persistent mode enables long-running operations to log intermediate accomplishments that remain visible:

```
   ✓ Built container
   ✓ Pushed to registry
   ✓ Updated load balancer
✓ Deploying application (took 45s)
```

**Behaviors:**
- StartPersistent begins the operation with a spinner
- LogAccomplishment prints an indented checkmark line (doesn't stop spinner)
- Each accomplishment remains visible as spinner continues
- FinishPersistent completes the operation and shows total time

**When to use:**
- Multi-step deployments or installations
- Batch processing with progress checkpoints
- Operations where users benefit from seeing incremental progress
- Situations where a simple spinner isn't informative enough

**Why persistent mode:**
- Provides visibility into long-running operations
- Shows progress without requiring user to watch constantly
- Creates a visible record of what was accomplished

## Pause/Resume Pattern

The pause/resume mechanism enables interactive prompts during progress operations:

1. **Pause** - Cancels spinners, clears terminal line, marks state as paused
2. **User Interaction** - Application shows prompt, gets input (no library involvement)
3. **Resume** - Restarts spinner for currently active operation

**Design considerations:**
- Library doesn't handle the actual user interaction
- Pause simply stops spinners and clears the line
- Resume restarts where it left off
- No state is lost during pause
- Nested operations maintain their hierarchy through pause/resume

**Why separate from library:**
- User input mechanisms vary (stdin, TUI libraries, etc.)
- Library focuses on progress display, not interaction
- Keeps API simple and focused

## Timing Display Strategy

Operation duration is displayed only when meaningful:

- **Threshold:** Only show duration for operations > 100ms
- **Precision:** Round to nearest 10ms for readability
- **Format:** "took 2s" or "took 450ms"

**Why this approach:**
- Sub-100ms operations are effectively instantaneous to users
- Prevents noise from trivial operations
- Rounding makes timing easier to read at a glance

## Terminal Control Strategy

The library manages cursor visibility and line clearing:

- **Hide cursor** when spinners are active (cleaner visual appearance)
- **Show cursor** when paused or after completion
- **Clear line** before writing to remove old spinner frames
- **Track cursor state** with atomic flag to avoid redundant control sequences

**Why manage cursor:**
- Spinner animation looks cleaner without visible cursor
- Line clearing prevents visual artifacts
- Proper cleanup ensures terminal is left in good state

## Zero Configuration Philosophy

The library requires no configuration or initialization beyond calling the constructor:

- Sensible defaults for all behavior (timing thresholds, spinner speed, etc.)
- No configuration files or environment variables
- No global state or singletons
- Works immediately with `NewProgressDisplay(os.Stdout)`

**Why zero configuration:**
- Reduces cognitive load for users
- Eliminates entire class of setup errors
- Makes library easier to adopt and use
- Fewer decisions for users to make

## Error Handling Model

Progress operations follow a simple error model:

- Most methods return error for consistency but rarely fail in practice
- Actual progress work (user's business logic) is separate from progress display
- Fail() method takes both message and error for context
- Display errors (output problems) don't stop the application

**Design principle:** Progress display failures should never break the application. If output fails, the operation should continue and succeed/fail based on actual work, not display.

## Noop Implementation

The library includes a no-op implementation primarily for testing purposes:

- Satisfies the same interface as the full implementation
- All methods succeed silently without output
- No goroutines, no terminal control, no overhead

**Use cases:**
- Unit testing without terminal output
- Silent mode in CI environments
- Disabling progress display when not needed

This is a supporting feature, not a core architectural pattern. The interface-based design makes it straightforward to provide alternative implementations when needed.

## Dependencies Strategy

Minimal external dependencies:
- **charmbracelet/huh/spinner** - Spinner animations
- **charmbracelet/lipgloss** - Terminal styling

**Why these dependencies:**
- huh/spinner provides smooth, well-tested spinner implementations
- lipgloss enables styled terminal output
- Both are stable, maintained libraries from the same ecosystem
- No need to implement terminal control from scratch

**Dependency philosophy:** Keep dependencies minimal since this is a library. New dependencies require strong justification (significant value add, stable, well-maintained).
