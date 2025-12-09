package dbutils

import (
	"errors"
	"fmt"

	"gorm.io/gorm"
)

func IsRecordNotFound(err error) bool {
	if err == nil {
		return false
	}
	return errors.Is(err, gorm.ErrRecordNotFound)
}

func IgnoreRecordNotFound(err error) error {
	if IsRecordNotFound(err) {
		return nil
	}
	return err
}

func WrapIfNotFound(err error, message string) error {
	if IsRecordNotFound(err) {
		return nil
	}
	return fmt.Errorf("%s: %w", message, err)
}
