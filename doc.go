// Package nesgress provides hierarchical progress reporting with nested spinner displays.
//
// The name "nesgress" is a pun on "nested progress", reflecting the library's ability
// to display hierarchical, nested progress operations in the terminal.
//
// # Features
//
//   - Hierarchical progress display with parent-child relationships
//   - Thread-safe concurrent operation support
//   - Spinner animations using charmbracelet/huh
//   - Success/failure indicators with timing information
//   - Persistent mode for long-running operations with accomplishments
//   - Pause/resume support for interactive prompts
//
// # Basic Usage
//
//	display := nesgress.NewProgressDisplay(os.Stdout)
//	defer display.Close()
//
//	display.Start("Installing packages")
//	// ... do work ...
//	display.Finish("Packages installed")
//
// # Nested Operations
//
//	display.Start("Setting up environment")
//	    display.Start("Downloading dependencies")
//	    display.Finish("Dependencies downloaded")
//
//	    display.Start("Compiling")
//	    display.Finish("Compiled")
//	display.Finish("Environment ready")
//
// # Persistent Mode
//
// For long-running operations where you want to show intermediate accomplishments:
//
//	display.StartPersistent("Deploying application")
//	display.LogAccomplishment("Built container")
//	display.LogAccomplishment("Pushed to registry")
//	display.FinishPersistent("Deployment complete")
//
// # Pause/Resume for Interactive Input
//
// When you need to prompt for user input:
//
//	display.Pause()   // Stops spinners, clears line
//	// ... show prompt, get input ...
//	display.Resume()  // Restarts spinners
//
// # Noop Implementation
//
// For testing or when progress display should be disabled:
//
//	display := nesgress.NewNoopProgressDisplay()
package nesgress
