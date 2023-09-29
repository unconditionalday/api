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
		t.Errorf("Errore inaspettato: %v", err)
	}
	if !result {
		t.Errorf("La versione v1 dovrebbe essere inferiore a v2")
	}

	// Test with v1 and v2 versions where v1 is greater
	v1 = "2023.03.01"
	v2 = "2023.02.01"
	result, err = calver.Lower(v1, v2)
	if err != nil {
		t.Errorf("Errore inaspettato: %v", err)
	}
	if result {
		t.Errorf("La versione v1 non dovrebbe essere inferiore a v2")
	}

	// Test with v1 and v2 versions where v1 is invalid
	v1 = "invalid"
	v2 = "2023.02.01"
	_, err = calver.Lower(v1, v2)
	if err == nil {
		t.Errorf("Dovrebbe esserci un errore per una versione non valida")
	}
}
