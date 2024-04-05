package config

import (
	"os"
	"testing"
)

type MockLoadEnv struct {
	Calls int
}

func (m *MockLoadEnv) Load(...string) error {
	m.Calls++
	os.Setenv("MONGO_USERNAME", "test")
	os.Setenv("MONGO_PASSWORD", "test")
	os.Setenv("MONGO_URL", "test")
	return nil
}

func TestGetEnvironment(t *testing.T) {
	mockLoadEnv := &MockLoadEnv{}
	GetEnvironment(mockLoadEnv.Load)
	if mockLoadEnv.Calls != 1 {
		t.Errorf("Expected mockLoadEnv to be called once, got %d", mockLoadEnv.Calls)
	}
	GetEnvironment(mockLoadEnv.Load)
	if mockLoadEnv.Calls != 1 {
		t.Errorf("Expected mockLoadEnv to be called once, got %d", mockLoadEnv.Calls)
	}

}
