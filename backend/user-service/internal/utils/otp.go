package utils

import (
	"crypto/rand"
	"fmt"
	"math/big"
)

func GenerateOTP(n int) (string, error) {
	if n <= 0 {
		return "", fmt.Errorf("OTP length must be positive")
	}

	max := big.NewInt(1)
	for i := 0; i < n; i++ {
		max.Mul(max, big.NewInt(10))
	}

	n64, err := rand.Int(rand.Reader, max)
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("%0*d", n, n64.Int64()), nil
}
