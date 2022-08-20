package gpooling

// MockedGPoolingImpl - mocking
type MockedGPoolingImpl struct {
}

// Release - release all gorotine
func (p *MockedGPoolingImpl) Release() {

}

// Running - returns the number of the currently running goroutines.
func (p *MockedGPoolingImpl) Running() int {
	return 0
}

// Submit - submit a task to this pool
func (p *MockedGPoolingImpl) Submit(task func()) {
	go task()
}
