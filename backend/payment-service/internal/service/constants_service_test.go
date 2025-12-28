package service

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewConstantsService(t *testing.T) {
	service := NewConstantsService()

	assert.NotNil(t, service)
	assert.IsType(t, &ConstantsServiceImpl{}, service)
}

func TestGetBanks_Success(t *testing.T) {
	service := NewConstantsService()

	ctx := context.Background()

	banks, err := service.GetBanks(ctx)

	assert.NoError(t, err)
	assert.NotNil(t, banks)
	assert.Greater(t, len(banks), 0, "Should return at least one bank")

	// Verify structure of first bank
	if len(banks) > 0 {
		firstBank := banks[0]
		assert.NotEmpty(t, firstBank.Code, "Bank code should not be empty")
		assert.NotEmpty(t, firstBank.ShortName, "Bank short name should not be empty")
		assert.NotEmpty(t, firstBank.Name, "Bank name should not be empty")
	}
}

func TestGetBanks_ReturnsConsistentData(t *testing.T) {
	service := NewConstantsService()

	ctx := context.Background()

	// Call multiple times to ensure consistency
	banks1, err1 := service.GetBanks(ctx)
	banks2, err2 := service.GetBanks(ctx)

	assert.NoError(t, err1)
	assert.NoError(t, err2)
	assert.Equal(t, len(banks1), len(banks2), "Should return same number of banks")

	// Verify known banks exist
	bankCodes := make(map[string]bool)
	for _, bank := range banks1 {
		bankCodes[bank.Code] = true
	}

	expectedBanks := []string{"VCB", "TCB", "MB", "BIDV", "ACB"}
	for _, code := range expectedBanks {
		assert.True(t, bankCodes[code], "Expected bank %s to be in the list", code)
	}
}
