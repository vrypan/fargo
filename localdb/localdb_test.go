package localdb

import (
	"testing"
)

func TestBasic(t *testing.T) {
	err := Open()
	if err != nil {
		t.Fatalf("Failed to open database: %v", err)
	}
	defer Close()

	key := "testKey"
	value := "testValue"

	// Store the key/value
	err = Set(key, []byte(value))
	if err != nil {
		t.Fatalf("Failed to store data: %v", err)
	}

	// Retrieve the key/value
	retrievedValue, err := Get(key)
	if err != nil {
		t.Fatalf("Failed to retrieve data: %v", err)
	}

	if string(retrievedValue) != value {
		t.Errorf("Expected value '%v', got '%v'", value, retrievedValue)
	}
}
