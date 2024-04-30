package domain

import (
	"errors"
	"fmt"
	"testing"
)

func TestError_WrappedErrApplicationGeneric_ExpectIsErrApplicationGeneric(t *testing.T) {
	// Given
	// err := fmt.Errorf("%w:%v", ErrApplicationGeneric, "Generic Test Error")
	err := errors.Join(ErrApplication, fmt.Errorf("Generic Test Error"))
	if !errors.Is(err, ErrApplication) {
		t.Errorf("Expected error to be of type ErrorApplicationGeneric, got %v", err)
	}
}

func TestError_ErrApplicationGeneric_ExpectIsErrApplicationGeneric(t *testing.T) {
	// Given
	err := ErrApplication
	if !errors.Is(err, ErrApplication) {
		t.Errorf("Expected error to be of type ErrorApplicationGeneric, got %v", err)
	}
}

func TestError_ErrSpecific_ExpectIsErrSpecific(t *testing.T) {
	// Given
	var ErrSpecific = fmt.Errorf("%w:%v", ErrApplication, "Specific Test Error")
	err := ErrSpecific
	if !errors.Is(err, ErrSpecific) {
		t.Errorf("Expected error to be of type ErrSpecific, got %v", err)
	}
}