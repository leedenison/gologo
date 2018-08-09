package gologo

// Lifecycle : stores the array of registered Lifecycle functions
type Lifecycle []func()

// New : Create a new lifecycle
func New() *Lifecycle {
	s := make(Lifecycle, 0)
	return &s
}

// Cleanup : Call all lifecycle functions for cleanup
func (lifecycle *Lifecycle) Cleanup() {
	for _, callback := range *lifecycle {
		callback()
	}
}

// RegisterCleanup : add the provided callback to the cleaup lifecycle
// function list
func (lifecycle *Lifecycle) RegisterCleanup(callback func()) {
	s := append(*lifecycle, callback)
	lifecycle = &s
}
