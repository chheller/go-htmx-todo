package server

import (
	"errors"
	"fmt"
);


func HellowWorld(name string) (string, error) {
	if name == "" {
		return "", errors.New("empty name")
	}
	return fmt.Sprintf("Hello, %v!", name), nil
}