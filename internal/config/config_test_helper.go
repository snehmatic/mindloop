//go:build test

package config

// ResetForTest resets the singleton for testing purposes.
func ResetForTest() {
	mu.Lock()
	defer mu.Unlock()
	instance = nil
	onChangeFn = nil

	userMu.Lock()
	userInstance = nil
	userMu.Unlock()
}
