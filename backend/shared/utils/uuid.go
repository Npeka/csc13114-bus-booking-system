package utils

import (
	"fmt"

	"github.com/google/uuid"
)

// ParseUUIDs converts a slice of UUID strings to a slice of uuid.UUID
// Returns error if any string is not a valid UUID
func ParseUUIDs(ids []string) ([]uuid.UUID, error) {
	result := make([]uuid.UUID, 0, len(ids))
	for _, idStr := range ids {
		id, err := uuid.Parse(idStr)
		if err != nil {
			return nil, fmt.Errorf("invalid UUID: %s", idStr)
		}
		result = append(result, id)
	}
	return result, nil
}

// ParseUUIDsWithContext is like ParseUUIDs but includes context in error message
func ParseUUIDsWithContext(ids []string, context string) ([]uuid.UUID, error) {
	result := make([]uuid.UUID, 0, len(ids))
	for _, idStr := range ids {
		id, err := uuid.Parse(idStr)
		if err != nil {
			return nil, fmt.Errorf("invalid %s: %s", context, idStr)
		}
		result = append(result, id)
	}
	return result, nil
}
