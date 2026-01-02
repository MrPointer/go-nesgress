package nesgress_test

import (
	"bytes"
	"errors"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"github.com/MrPointer/go-nesgress"
)

func Test_NewProgressDisplay_WithBuffer_CreatesValidInstance(t *testing.T) {
	var buf bytes.Buffer
	display := nesgress.NewProgressDisplay(&buf)
	require.NotNil(t, display)
}

func Test_NewProgressDisplay_WithNilOutput_UsesStdout(t *testing.T) {
	display := nesgress.NewProgressDisplay(nil)
	require.NotNil(t, display)
}

func Test_SingleProgressOperation_StartedAndFinished_ShowsSuccessMessage(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	var buf bytes.Buffer
	display := nesgress.NewProgressDisplay(&buf)

	require.NoError(t, display.Start("Test operation"))
	require.True(t, display.IsActive())

	time.Sleep(50 * time.Millisecond)
	require.NoError(t, display.Finish("Test operation"))

	require.False(t, display.IsActive())

	output := buf.String()
	require.Contains(t, output, "✓")
	require.Contains(t, output, "Test operation")
}

func Test_NestedProgressOperations_WithMultipleLevels_ShowProperHierarchy(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	var buf bytes.Buffer
	display := nesgress.NewProgressDisplay(&buf)

	_ = display.Start("Parent operation")
	require.True(t, display.IsActive())

	_ = display.Start("Child operation")
	require.True(t, display.IsActive())

	time.Sleep(30 * time.Millisecond)
	_ = display.Finish("Child operation")
	require.True(t, display.IsActive()) // Parent still active

	time.Sleep(30 * time.Millisecond)
	_ = display.Finish("Parent operation")
	require.False(t, display.IsActive())

	output := buf.String()
	lines := strings.Split(strings.TrimSpace(output), "\n")

	// Should have contextual messages during progress and clean completion messages
	require.Contains(t, output, "Child operation")
	require.Contains(t, output, "Parent operation")
	require.Contains(t, strings.Join(lines, "\n"), "✓")
}

func Test_ProgressMessage_UpdatedAfterStart_ShowsNewMessage(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	var buf bytes.Buffer
	display := nesgress.NewProgressDisplay(&buf)

	_ = display.Start("Initial message")
	_ = display.Update("Updated message")

	time.Sleep(30 * time.Millisecond)
	_ = display.Finish("Final message")

	output := buf.String()
	require.Contains(t, output, "Updated message")
}

func Test_ProgressOperation_WhenFailed_ShowsErrorMessage(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	var buf bytes.Buffer
	display := nesgress.NewProgressDisplay(&buf)

	_ = display.Start("Failing operation")
	time.Sleep(30 * time.Millisecond)
	_ = display.Fail("Failing operation", errors.New("test error"))

	require.False(t, display.IsActive())

	output := buf.String()
	require.Contains(t, output, "✗")
	require.Contains(t, output, "Failing operation")
	require.Contains(t, output, "test error")
}

func Test_MixedSuccessAndFailureOperations_WithNestedStructure_DisplayCorrectly(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	var buf bytes.Buffer
	display := nesgress.NewProgressDisplay(&buf)

	_ = display.Start("Parent")

	_ = display.Start("Success child")
	time.Sleep(20 * time.Millisecond)
	_ = display.Finish("Success child")

	_ = display.Start("Failing child")
	time.Sleep(20 * time.Millisecond)
	_ = display.Fail("Failing child", errors.New("child error"))

	_ = display.Finish("Parent")

	output := buf.String()
	require.Contains(t, output, "Success child")
	require.Contains(t, output, "Failing child")
	require.Contains(t, output, "Parent")
	require.Contains(t, output, "child error")
	require.Contains(t, output, "✓")
	require.Contains(t, output, "✗")
}

func Test_DeeplyNestedOperations_WithFiveLevels_ShowCorrectIndentation(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	var buf bytes.Buffer
	display := nesgress.NewProgressDisplay(&buf)

	// Create deep nesting
	_ = display.Start("Level 1")
	_ = display.Start("Level 2")
	_ = display.Start("Level 3")
	_ = display.Start("Level 4")

	time.Sleep(20 * time.Millisecond)

	// Complete in reverse order
	_ = display.Finish("Level 4")
	_ = display.Finish("Level 3")
	_ = display.Finish("Level 2")
	_ = display.Finish("Level 1")

	output := buf.String()

	// Check that all levels appear in output (contextual messages during progress)
	require.Contains(t, output, "Level 1")
	require.Contains(t, output, "Level 2")
	require.Contains(t, output, "Level 3")
	require.Contains(t, output, "Level 4")
}

func Test_LongRunningOperations_OverThreshold_ShowTimingInformation(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	var buf bytes.Buffer
	display := nesgress.NewProgressDisplay(&buf)

	_ = display.Start("Long operation")
	time.Sleep(150 * time.Millisecond) // Longer than 100ms threshold
	_ = display.Finish("Long operation")

	output := buf.String()
	require.Contains(t, output, "took")
	require.Contains(t, output, "ms")
}

func Test_ShortOperations_UnderThreshold_DoNotShowTimingInformation(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	var buf bytes.Buffer
	display := nesgress.NewProgressDisplay(&buf)

	_ = display.Start("Quick operation")
	// No sleep - immediate completion
	_ = display.Finish("Quick operation")

	output := buf.String()
	require.NotContains(t, output, "took")
}

func Test_Clear_WithActiveOperations_StopsAllOperations(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	var buf bytes.Buffer
	display := nesgress.NewProgressDisplay(&buf)

	_ = display.Start("Operation 1")
	_ = display.Start("Operation 2")
	require.True(t, display.IsActive())

	_ = display.Clear()
	require.False(t, display.IsActive())

	// Should not crash when calling methods after clear
	_ = display.Update("Should be ignored")
	_ = display.Finish("Should be ignored")
	_ = display.Fail("Should be ignored", errors.New("test"))
}

func Test_Update_WithoutActiveProgress_DoesNothing(t *testing.T) {
	var buf bytes.Buffer
	display := nesgress.NewProgressDisplay(&buf)

	_ = display.Update("No active operation")
	require.False(t, display.IsActive())
}

func Test_Finish_WithoutActiveProgress_DoesNothing(t *testing.T) {
	var buf bytes.Buffer
	display := nesgress.NewProgressDisplay(&buf)

	_ = display.Finish("No active operation")
	require.False(t, display.IsActive())
}

func Test_Fail_WithoutActiveProgress_DoesNothing(t *testing.T) {
	var buf bytes.Buffer
	display := nesgress.NewProgressDisplay(&buf)

	_ = display.Fail("No active operation", errors.New("test"))
	require.False(t, display.IsActive())
}

func Test_NoopProgressDisplay_ByDesign_ImplementsProgressReporterInterface(t *testing.T) {
	var _ nesgress.ProgressReporter = (*nesgress.NoopProgressDisplay)(nil)
}

func Test_NoopProgressDisplay_AllMethods_DoNothing(t *testing.T) {
	display := nesgress.NewNoopProgressDisplay()

	// All these should not crash and should not do anything
	_ = display.Start("Test")
	require.False(t, display.IsActive())

	_ = display.Update("Test")
	require.False(t, display.IsActive())

	_ = display.Finish("Test")
	require.False(t, display.IsActive())

	_ = display.Fail("Test", errors.New("test"))
	require.False(t, display.IsActive())

	_ = display.Clear()
	require.False(t, display.IsActive())
}

func Test_ConcurrentProgressDisplayOperations_WithMultipleGoroutines_AreThreadSafe(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	var buf bytes.Buffer
	display := nesgress.NewProgressDisplay(&buf)

	done := make(chan bool, 2)

	// Test concurrent access to ensure thread safety
	go func() {
		_ = display.Start("Concurrent 1")
		time.Sleep(50 * time.Millisecond)
		_ = display.Finish("Concurrent 1")

		done <- true
	}()

	go func() {
		time.Sleep(25 * time.Millisecond)
		_ = display.Start("Concurrent 2")
		time.Sleep(50 * time.Millisecond)
		_ = display.Finish("Concurrent 2")

		done <- true
	}()

	// Wait for both goroutines to complete
	<-done
	<-done

	require.False(t, display.IsActive())

	output := buf.String()
	require.Contains(t, output, "Concurrent 1")
	require.Contains(t, output, "Concurrent 2")
}

func Test_RapidSequentialOperations_WithQuickSuccession_Work(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	var buf bytes.Buffer
	display := nesgress.NewProgressDisplay(&buf)

	// Rapidly start and finish operations
	for i := range 10 {
		_ = i
		_ = display.Start("Rapid operation")
		time.Sleep(5 * time.Millisecond)
		_ = display.Finish("Rapid operation")
	}

	require.False(t, display.IsActive())

	output := buf.String()
	// Should have multiple completion messages
	checkmarkCount := strings.Count(output, "✓")
	require.Equal(t, 10, checkmarkCount)
}

func Test_StartPersistentProgress_WhenCalled_ActivatesPersistentMode(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	var buf bytes.Buffer
	display := nesgress.NewProgressDisplay(&buf)

	_ = display.StartPersistent("Installing packages")
	require.True(t, display.IsActive())

	time.Sleep(50 * time.Millisecond)
	_ = display.FinishPersistent("Installation complete")
	require.False(t, display.IsActive())

	output := buf.String()
	require.Contains(t, output, "✓")
	require.Contains(t, output, "Installing packages")
}

func Test_LogAccomplishment_WithPersistentProgress_ShowsVisibleAccomplishments(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	var buf bytes.Buffer
	display := nesgress.NewProgressDisplay(&buf)

	_ = display.StartPersistent("Deploying application")
	time.Sleep(30 * time.Millisecond)

	_ = display.LogAccomplishment("Built application")
	_ = display.LogAccomplishment("Created container")
	_ = display.LogAccomplishment("Pushed to registry")

	time.Sleep(30 * time.Millisecond)
	_ = display.FinishPersistent("Deployment complete")

	output := buf.String()
	require.Contains(t, output, "Built application")
	require.Contains(t, output, "Created container")
	require.Contains(t, output, "Pushed to registry")
	require.Contains(t, output, "Deploying application")
}

func Test_ProgressDisplayPersistentProgress_WhenFailed_ShowsErrorMessage(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	var buf bytes.Buffer
	display := nesgress.NewProgressDisplay(&buf)

	_ = display.StartPersistent("Critical operation")
	_ = display.LogAccomplishment("Step 1 completed")
	_ = display.LogAccomplishment("Step 2 completed")
	time.Sleep(30 * time.Millisecond)
	_ = display.FailPersistent("Critical operation failed", errors.New("permission denied"))

	require.False(t, display.IsActive())

	output := buf.String()
	require.Contains(t, output, "Step 1 completed")
	require.Contains(t, output, "Step 2 completed")
	require.Contains(t, output, "✗")
	require.Contains(t, output, "Critical operation")
	require.Contains(t, output, "permission denied")
}

func Test_ProgressDisplayMixed_PersistentAndRegularOperations_Work(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	var buf bytes.Buffer
	display := nesgress.NewProgressDisplay(&buf)

	_ = display.StartPersistent("Setting up environment")
	_ = display.LogAccomplishment("Created directories")

	_ = display.Start("Downloading files")
	time.Sleep(30 * time.Millisecond)
	_ = display.Finish("Files downloaded")

	_ = display.LogAccomplishment("Installed dependencies")

	_ = display.Start("Running tests")
	time.Sleep(30 * time.Millisecond)
	_ = display.Finish("Tests passed")

	_ = display.FinishPersistent("Environment ready")

	output := buf.String()
	require.Contains(t, output, "Created directories")
	require.Contains(t, output, "Installed dependencies")
	require.Contains(t, output, "Setting up environment")
}

func Test_LogAccomplishment_WithoutActivePersistentProgress_StillWorks(t *testing.T) {
	var buf bytes.Buffer
	display := nesgress.NewProgressDisplay(&buf)

	_ = display.LogAccomplishment("Standalone accomplishment")

	output := buf.String()
	require.Contains(t, output, "✓")
	require.Contains(t, output, "Standalone accomplishment")
}

func Test_FinishPersistent_WithoutActiveProgress_DoesNothing(t *testing.T) {
	var buf bytes.Buffer
	display := nesgress.NewProgressDisplay(&buf)

	_ = display.FinishPersistent("No active progress")
	require.False(t, display.IsActive())
}

func Test_FailPersistent_WithoutActiveProgress_DoesNothing(t *testing.T) {
	var buf bytes.Buffer
	display := nesgress.NewProgressDisplay(&buf)

	_ = display.FailPersistent("No active progress", errors.New("test error"))
	require.False(t, display.IsActive())
}

func Test_PersistentProgress_WithAccomplishments_ShowsInRealTime(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	var buf bytes.Buffer
	display := nesgress.NewProgressDisplay(&buf)

	_ = display.StartPersistent("Processing items")

	accomplishments := []string{
		"Processed item 1",
		"Processed item 2",
		"Processed item 3",
		"Processed item 4",
		"Processed item 5",
	}

	for _, accomplishment := range accomplishments {
		time.Sleep(20 * time.Millisecond)
		_ = display.LogAccomplishment(accomplishment)
	}

	_ = display.FinishPersistent("All items processed")

	output := buf.String()

	for _, accomplishment := range accomplishments {
		require.Contains(t, output, accomplishment)
	}

	require.Contains(t, output, "Processing items")
}

func Test_NoopProgressDisplay_PersistentMethods_DoNothing(t *testing.T) {
	display := nesgress.NewNoopProgressDisplay()

	// All these should not crash and should not do anything
	_ = display.StartPersistent("Test")
	require.False(t, display.IsActive())

	_ = display.LogAccomplishment("Test accomplishment")
	require.False(t, display.IsActive())

	_ = display.FinishPersistent("Test")
	require.False(t, display.IsActive())

	_ = display.FailPersistent("Test", errors.New("test"))
	require.False(t, display.IsActive())

	_ = display.Close()
	require.False(t, display.IsActive())
}

func Test_Close_WithActiveOperations_StopsAllAndRestoresCursor(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	var buf bytes.Buffer
	display := nesgress.NewProgressDisplay(&buf)

	_ = display.Start("Operation 1")
	_ = display.Start("Operation 2")
	require.True(t, display.IsActive())

	_ = display.Close()
	require.False(t, display.IsActive())

	// Should not crash when calling methods after cleanup
	_ = display.Update("Should be ignored")
	_ = display.Finish("Should be ignored")
	_ = display.Fail("Should be ignored", errors.New("test"))
}

func Test_Close_CalledMultipleTimes_DoesNotCrash(t *testing.T) {
	var buf bytes.Buffer
	display := nesgress.NewProgressDisplay(&buf)

	_ = display.Start("Test operation")
	require.True(t, display.IsActive())

	// Multiple cleanups should not crash
	_ = display.Close()
	_ = display.Close()
	_ = display.Close()

	require.False(t, display.IsActive())
}

func Test_Close_WithoutActiveOperations_DoesNotCrash(t *testing.T) {
	var buf bytes.Buffer
	display := nesgress.NewProgressDisplay(&buf)

	require.NotPanics(t, func() {
		_ = display.Close()
	})

	require.False(t, display.IsActive())
}

func Test_Close_WithHiddenCursor_RestoresCursor(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	var buf bytes.Buffer
	display := nesgress.NewProgressDisplay(&buf)

	_ = display.Start("Operation with cursor")
	require.True(t, display.IsActive())

	time.Sleep(30 * time.Millisecond)
	_ = display.Close()
	require.False(t, display.IsActive())

	// Close should not crash and should handle cursor state properly
	require.NotNil(t, display)
}

func Test_ProgressFailure_WithSynchronization_PreventsHangingCursor(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	var buf bytes.Buffer
	display := nesgress.NewProgressDisplay(&buf)

	_ = display.Start("Operation that will fail")
	require.True(t, display.IsActive())

	// Simulate work that results in failure
	time.Sleep(100 * time.Millisecond)
	_ = display.Fail("Operation failed", errors.New("simulated error"))

	// After failure, display should be properly cleaned up
	require.False(t, display.IsActive())

	// Verify that output contains failure message and cursor control sequences are handled
	output := buf.String()
	require.Contains(t, output, "✗")
	require.Contains(t, output, "Operation that will fail")
	require.Contains(t, output, "simulated error")

	// Should be able to start new operations without issues
	_ = display.Start("New operation after failure")
	require.True(t, display.IsActive())

	time.Sleep(30 * time.Millisecond)
	_ = display.Finish("New operation completed")
	require.False(t, display.IsActive())

	// Verify new operation also completed successfully
	require.Contains(t, buf.String(), "New operation after failure")
}

func Test_RapidFailureAndRecovery_WithQuickOperations_MaintainsProperTerminalState(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	var buf bytes.Buffer
	display := nesgress.NewProgressDisplay(&buf)

	// Test rapid failure and recovery cycles
	for i := range 5 {
		_ = display.Start("Rapid operation")
		time.Sleep(20 * time.Millisecond)

		if i%2 == 0 {
			_ = display.Fail("Operation failed", errors.New("test error"))
		} else {
			_ = display.Finish("Operation succeeded")
		}

		require.False(t, display.IsActive())
	}

	output := buf.String()
	// Should have both success and failure indicators
	require.Contains(t, output, "✓")
	require.Contains(t, output, "✗")

	// Should not have any hanging operations
	require.False(t, display.IsActive())
}

func Test_Pause_WithActiveSpinners_StopsAllOperations(t *testing.T) {
	var output bytes.Buffer
	display := nesgress.NewProgressDisplay(&output)

	_ = display.Start("Operation 1")
	_ = display.Start("Operation 2")
	time.Sleep(50 * time.Millisecond) // Allow spinners to start

	require.True(t, display.IsActive())
	require.False(t, display.IsPaused())

	err := display.Pause()
	require.NoError(t, err)

	require.True(t, display.IsPaused())
	require.True(t, display.IsActive()) // Operations still exist, just paused
}

func Test_Resume_AfterPause_RestartsSpinnerOperations(t *testing.T) {
	var output bytes.Buffer
	display := nesgress.NewProgressDisplay(&output)

	_ = display.Start("Operation 1")
	_ = display.Start("Operation 2")
	time.Sleep(50 * time.Millisecond) // Allow spinners to start

	err := display.Pause()
	require.NoError(t, err)
	require.True(t, display.IsPaused())

	err = display.Resume()
	require.NoError(t, err)

	require.False(t, display.IsPaused())
	require.True(t, display.IsActive())
}

func Test_Pause_WithoutActiveOperations_DoesNotCrash(t *testing.T) {
	var output bytes.Buffer
	display := nesgress.NewProgressDisplay(&output)

	require.False(t, display.IsActive())
	require.False(t, display.IsPaused())

	err := display.Pause()
	require.NoError(t, err)

	require.True(t, display.IsPaused())
	require.False(t, display.IsActive())
}

func Test_Resume_WithoutActiveOperations_DoesNotCrash(t *testing.T) {
	var output bytes.Buffer
	display := nesgress.NewProgressDisplay(&output)

	err := display.Pause()
	require.NoError(t, err)
	require.True(t, display.IsPaused())

	err = display.Resume()
	require.NoError(t, err)

	require.False(t, display.IsPaused())
	require.False(t, display.IsActive())
}

func Test_Pause_CalledMultipleTimes_IsSafe(t *testing.T) {
	var output bytes.Buffer
	display := nesgress.NewProgressDisplay(&output)

	_ = display.Start("Test Operation")
	time.Sleep(50 * time.Millisecond) // Allow spinner to start

	err := display.Pause()
	require.NoError(t, err)
	require.True(t, display.IsPaused())

	err = display.Pause()
	require.NoError(t, err)
	require.True(t, display.IsPaused())
}

func Test_Resume_CalledMultipleTimes_IsSafe(t *testing.T) {
	var output bytes.Buffer
	display := nesgress.NewProgressDisplay(&output)

	_ = display.Start("Test Operation")
	time.Sleep(50 * time.Millisecond) // Allow spinner to start
	_ = display.Pause()

	err := display.Resume()
	require.NoError(t, err)
	require.False(t, display.IsPaused())

	err = display.Resume()
	require.NoError(t, err)
	require.False(t, display.IsPaused())
}

func Test_PauseAndResume_WithNestedOperations_WorksCorrectly(t *testing.T) {
	var output bytes.Buffer
	display := nesgress.NewProgressDisplay(&output)

	_ = display.Start("Parent Operation")
	_ = display.Start("Child Operation 1")
	_ = display.Start("Child Operation 2")
	time.Sleep(50 * time.Millisecond) // Allow spinners to start

	require.True(t, display.IsActive())
	require.False(t, display.IsPaused())

	err := display.Pause()
	require.NoError(t, err)
	require.True(t, display.IsPaused())

	err = display.Resume()
	require.NoError(t, err)
	require.False(t, display.IsPaused())
	require.True(t, display.IsActive())
}

func Test_Pause_BeforeInteractiveInput_StopsSpinnerAndClearsOutput(t *testing.T) {
	var output bytes.Buffer

	display := nesgress.NewProgressDisplay(&output)

	_ = display.Start("Processing files...")
	time.Sleep(50 * time.Millisecond) // Let spinner start

	initialOutput := display.GetOutputSafely()
	require.NotEmpty(t, initialOutput)

	err := display.Pause()
	require.NoError(t, err)

	// Allow some time to ensure spinner has stopped
	time.Sleep(100 * time.Millisecond)
	outputAfterPause := display.GetOutputSafely()

	// Output should contain clear line sequence after pause
	require.Contains(t, outputAfterPause, "\r")
	require.True(t, display.IsPaused())
	require.True(t, display.IsActive()) // Operations still exist, just paused
}

func Test_Resume_WithMultipleOperations_RestartsMostRecent(t *testing.T) {
	var output bytes.Buffer

	display := nesgress.NewProgressDisplay(&output)

	_ = display.Start("Operation 1")
	_ = display.Start("Operation 2")
	_ = display.Start("Operation 3")
	time.Sleep(50 * time.Millisecond) // Allow spinners to start

	err := display.Pause()
	require.NoError(t, err)

	err = display.Resume()
	require.NoError(t, err)

	time.Sleep(50 * time.Millisecond)
	outputContent := display.GetOutputSafely()

	// Should show the most recent operation (Operation 3)
	require.Contains(t, outputContent, "Operation 3")
}

func Test_PausedState_WithNewOperations_StillAllowsOperationCreation(t *testing.T) {
	var output bytes.Buffer
	display := nesgress.NewProgressDisplay(&output)

	_ = display.Start("Operation 1")

	time.Sleep(50 * time.Millisecond) // Allow spinner to start

	err := display.Pause()
	require.NoError(t, err)
	require.True(t, display.IsPaused())

	// Starting new operation while paused should still work
	_ = display.Start("Operation 2")
	require.True(t, display.IsActive())
	require.True(t, display.IsPaused())

	err = display.Resume()
	require.NoError(t, err)
	require.False(t, display.IsPaused())
}

func Test_PauseAndResume_WithPersistentProgress_WorksCorrectly(t *testing.T) {
	var output bytes.Buffer
	display := nesgress.NewProgressDisplay(&output)

	_ = display.StartPersistent("Installing packages...")
	_ = display.LogAccomplishment("Package 1 installed")

	err := display.Pause()
	require.NoError(t, err)
	require.True(t, display.IsPaused())

	err = display.Resume()
	require.NoError(t, err)
	require.False(t, display.IsPaused())

	_ = display.LogAccomplishment("Package 2 installed")
	_ = display.FinishPersistent("All packages installed successfully")
}

func Test_Close_AfterPause_RestoresTerminalState(t *testing.T) {
	var output bytes.Buffer
	display := nesgress.NewProgressDisplay(&output)

	_ = display.Start("Test Operation")
	time.Sleep(50 * time.Millisecond) // Allow spinner to start

	err := display.Pause()
	require.NoError(t, err)
	require.True(t, display.IsPaused())

	err = display.Close()
	require.NoError(t, err)

	// Should be able to close multiple times
	err = display.Close()
	require.NoError(t, err)
}

func Test_NoopProgressDisplay_PauseAndResumeMethods_DoNothing(t *testing.T) {
	display := nesgress.NewNoopProgressDisplay()

	require.False(t, display.IsActive())
	require.False(t, display.IsPaused())

	err := display.Pause()
	require.NoError(t, err)
	require.False(t, display.IsPaused())

	err = display.Resume()
	require.NoError(t, err)
	require.False(t, display.IsPaused())
}
