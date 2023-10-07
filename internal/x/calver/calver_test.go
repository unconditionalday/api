package calver_test

import (
	"testing"

	"github.com/unconditionalday/server/internal/x/calver"
)

func TestCalVerLower(t *testing.T) {
	calver := calver.New()

	// Testing with v1 and v2 versions where v1 is lower
	v1 := "2023.01.01"
	v2 := "2023.02.01"
	result, err := calver.Lower(v1, v2)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	if !result {
		t.Errorf("The v1 should be lower than v2")
	}

	// Test with v1 and v2 versions where v1 is greater
	v1 = "2023.03.01"
	v2 = "2023.02.01"
	result, err = calver.Lower(v1, v2)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	if result {
		t.Errorf("The v1 shouldn't be lower than v2")
	}

	// Test with v1 and v2 versions where v1 is invalid
	v1 = "invalid"
	v2 = "2023.02.01"
	_, err = calver.Lower(v1, v2)
	if err == nil {
		t.Errorf("Error expected")
	}
}
