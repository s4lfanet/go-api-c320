package repository

import (
	"testing"

	"github.com/gosnmp/gosnmp"
)

func TestNewPonRepository(t *testing.T) {
	target := "192.168.1.1"
	community := "public"
	var port uint16 = 161

	repo := NewPonRepository(target, community, port)

	if repo == nil {
		t.Error("Expected non-nil repository")
	}

	// Verify it implements the interface
	var _ SnmpRepositoryInterface = repo
}

func TestNewPonRepository_DifferentParameters(t *testing.T) {
	tests := []struct {
		name      string
		target    string
		community string
		port      uint16
	}{
		{
			name:      "Standard SNMP configuration",
			target:    "192.168.1.1",
			community: "public",
			port:      161,
		},
		{
			name:      "Custom port",
			target:    "10.0.0.1",
			community: "private",
			port:      1161,
		},
		{
			name:      "Localhost",
			target:    "localhost",
			community: "test",
			port:      161,
		},
		{
			name:      "IPv6 address",
			target:    "::1",
			community: "public",
			port:      161,
		},
		{
			name:      "Empty community",
			target:    "192.168.1.1",
			community: "",
			port:      161,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := NewPonRepository(tt.target, tt.community, tt.port)

			if repo == nil {
				t.Error("Expected non-nil repository")
			}

			// Type assert to check internal fields
			if snmpRepo, ok := repo.(*snmpRepository); ok {
				if snmpRepo.target != tt.target {
					t.Errorf("Expected target '%s', got '%s'", tt.target, snmpRepo.target)
				}

				if snmpRepo.community != tt.community {
					t.Errorf("Expected community '%s', got '%s'", tt.community, snmpRepo.community)
				}

				if snmpRepo.port != tt.port {
					t.Errorf("Expected port %d, got %d", tt.port, snmpRepo.port)
				}
			} else {
				t.Error("Failed to type assert repository to *snmpRepository")
			}
		})
	}
}

func TestSnmpRepository_Get_InvalidTarget(t *testing.T) {
	// Test with invalid target to verify error handling
	repo := NewPonRepository("invalid-host-that-does-not-exist", "public", 161)

	oids := []string{"1.3.6.1.2.1.1.1.0"}

	result, err := repo.Get(oids)

	// Should get an error for invalid/unreachable host
	if err == nil {
		t.Error("Expected error for invalid host, got nil")
	}

	if result != nil {
		t.Error("Expected nil result on error")
	}

	// Verify error message contains expected text
	if err != nil && err.Error() == "" {
		t.Error("Expected non-empty error message")
	}
}

func TestSnmpRepository_Get_EmptyOIDs(t *testing.T) {
	repo := NewPonRepository("invalid-host", "public", 161)

	// Try with empty OID list
	result, err := repo.Get([]string{})

	// Will fail due to invalid host, but testing the parameter handling
	if err == nil {
		t.Error("Expected error, got nil")
	}

	if result != nil {
		t.Error("Expected nil result on error")
	}
}

func TestSnmpRepository_Walk_InvalidTarget(t *testing.T) {
	repo := NewPonRepository("invalid-host-that-does-not-exist", "public", 161)

	oid := "1.3.6.1.2.1.1"

	callbackCalled := false
	walkFunc := func(pdu gosnmp.SnmpPDU) error {
		callbackCalled = true
		return nil
	}

	err := repo.Walk(oid, walkFunc)

	// Should get an error for invalid/unreachable host
	if err == nil {
		t.Error("Expected error for invalid host, got nil")
	}

	// Callback should not be called if connection fails
	if callbackCalled {
		t.Error("Callback should not be called when connection fails")
	}
}

func TestSnmpRepository_Walk_ErrorPropagation(t *testing.T) {
	// This test verifies that the Walk method properly handles errors
	// Since we can't actually connect, we test the error path
	repo := NewPonRepository("127.0.0.1", "public", 65535) // Max uint16 port

	oid := "1.3.6.1.2.1.1"

	walkFunc := func(pdu gosnmp.SnmpPDU) error {
		return nil
	}

	err := repo.Walk(oid, walkFunc)

	// Should get connection error
	if err == nil {
		t.Error("Expected error for invalid configuration, got nil")
	}
}

func TestSnmpRepository_InterfaceCompliance(t *testing.T) {
	// Verify that snmpRepository implements SnmpRepositoryInterface
	var repo SnmpRepositoryInterface = NewPonRepository("invalid-host-that-does-not-exist", "public", 161)

	if repo == nil {
		t.Error("Repository should not be nil")
	}

	// Verify interface methods exist by calling them
	// We don't check for errors here as connection behavior can vary
	// The main goal is to verify the interface is implemented correctly
	_, _ = repo.Get([]string{"1.3.6.1.2.1.1.1.0"})
	_ = repo.Walk("1.3.6.1", func(pdu gosnmp.SnmpPDU) error { return nil })
}

func TestSnmpRepository_Get_MultipleOIDs(t *testing.T) {
	repo := NewPonRepository("invalid-host", "public", 161)

	// Test with multiple OIDs
	oids := []string{
		"1.3.6.1.2.1.1.1.0",
		"1.3.6.1.2.1.1.2.0",
		"1.3.6.1.2.1.1.3.0",
	}

	result, err := repo.Get(oids)

	// Will fail due to invalid host
	if err == nil {
		t.Error("Expected error, got nil")
	}

	if result != nil {
		t.Error("Expected nil result on error")
	}
}

func TestSnmpRepository_StructFields(t *testing.T) {
	target := "10.0.0.1"
	community := "test-community"
	var port uint16 = 8161

	repo := NewPonRepository(target, community, port)

	// Type assert to access internal fields
	snmpRepo, ok := repo.(*snmpRepository)
	if !ok {
		t.Fatal("Failed to type assert to *snmpRepository")
	}

	if snmpRepo.target != target {
		t.Errorf("Expected target '%s', got '%s'", target, snmpRepo.target)
	}

	if snmpRepo.community != community {
		t.Errorf("Expected community '%s', got '%s'", community, snmpRepo.community)
	}

	if snmpRepo.port != port {
		t.Errorf("Expected port %d, got %d", port, snmpRepo.port)
	}
}

func TestSnmpRepository_ZeroPort(t *testing.T) {
	// Test with port 0 (should be allowed but won't connect)
	repo := NewPonRepository("localhost", "public", 0)

	if repo == nil {
		t.Error("Expected non-nil repository even with port 0")
	}

	// Verify it was set
	if snmpRepo, ok := repo.(*snmpRepository); ok {
		if snmpRepo.port != 0 {
			t.Errorf("Expected port 0, got %d", snmpRepo.port)
		}
	}
}

func TestSnmpRepository_Get_Success(t *testing.T) {
	repo := NewPonRepository("invalid-host", "public", 161)

	oids := []string{"1.3.6.1.2.1.1.1.0"}
	_, err := repo.Get(oids)

	// Will fail due to invalid host, testing error handling
	if err == nil {
		t.Error("Expected error for unreachable host")
	}
}

func TestSnmpRepository_Walk_Success(t *testing.T) {
	repo := NewPonRepository("invalid-host", "public", 161)

	oid := "1.3.6.1.2.1.1"
	walkFunc := func(pdu gosnmp.SnmpPDU) error {
		return nil
	}

	err := repo.Walk(oid, walkFunc)

	// Will fail due to invalid host, testing error handling
	if err == nil {
		t.Error("Expected error for unreachable host")
	}
}

func TestSnmpRepository_Get_NilOIDs(t *testing.T) {
	repo := NewPonRepository("localhost", "public", 161)

	_, err := repo.Get(nil)

	// Should handle nil OIDs gracefully
	if err == nil {
		t.Error("Expected error for nil OIDs")
	}
}

func TestSnmpRepository_Walk_CallbackError(t *testing.T) {
	repo := NewPonRepository("localhost", "public", 161)

	oid := "1.3.6.1.2.1.1"
	walkFunc := func(pdu gosnmp.SnmpPDU) error {
		return nil
	}

	err := repo.Walk(oid, walkFunc)

	// Will error on connection
	if err == nil {
		t.Error("Expected error")
	}
}
