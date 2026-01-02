package nesgress

// NoopProgressDisplay is a progress display that does nothing.
type NoopProgressDisplay struct{}

var _ ProgressReporter = (*NoopProgressDisplay)(nil)

// NewNoopProgressDisplay creates a progress display that does nothing.
func NewNoopProgressDisplay() *NoopProgressDisplay {
	return &NoopProgressDisplay{}
}

// Start does nothing.
func (n *NoopProgressDisplay) Start(message string) error {
	return nil
}

// Update does nothing.
func (n *NoopProgressDisplay) Update(message string) error {
	return nil
}

// Finish does nothing.
func (n *NoopProgressDisplay) Finish(message string) error {
	return nil
}

// Fail does nothing.
func (n *NoopProgressDisplay) Fail(message string, err error) error {
	return nil
}

// IsActive always returns false.
func (n *NoopProgressDisplay) IsActive() bool { return false }

// Clear does nothing.
func (n *NoopProgressDisplay) Clear() error { return nil }

// Pause does nothing.
func (n *NoopProgressDisplay) Pause() error {
	return nil
}

// Resume does nothing.
func (n *NoopProgressDisplay) Resume() error {
	return nil
}

// IsPaused always returns false.
func (n *NoopProgressDisplay) IsPaused() bool { return false }

// StartPersistent does nothing.
func (n *NoopProgressDisplay) StartPersistent(message string) error {
	return nil
}

// LogAccomplishment does nothing.
func (n *NoopProgressDisplay) LogAccomplishment(message string) error {
	return nil
}

// FinishPersistent does nothing.
func (n *NoopProgressDisplay) FinishPersistent(message string) error {
	return nil
}

// FailPersistent does nothing.
func (n *NoopProgressDisplay) FailPersistent(message string, err error) error {
	return nil
}

// Close does nothing.
func (n *NoopProgressDisplay) Close() error { return nil }
