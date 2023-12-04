package hash

import (
	"github.com/stretchr/testify/require"
	"testing"
)

func TestGenerateAndCompare(t *testing.T) {
	password := "valentin228Tubik@s"
	hashedPassword, err := HashPassword(password)
	require.NoError(t, err)

	inputPassword := "valentin228Tubik@s"
	require.True(t, ComparePassword(hashedPassword, inputPassword))
}
