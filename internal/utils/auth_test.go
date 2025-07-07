package utils

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestHashPassword(t *testing.T) {
	password := "testpassword123"

	hash, err := HashPassword(password)
	require.NoError(t, err)
	assert.NotEmpty(t, hash)

	// Should start with $argon2id$
	assert.True(t, strings.HasPrefix(hash, "$argon2id$"))

	// Should have 6 parts when split by $
	parts := strings.Split(hash, "$")
	assert.Len(t, parts, 6)
}

func TestVerifyPassword(t *testing.T) {
	password := "testpassword123"
	wrongPassword := "wrongpassword"

	hash, err := HashPassword(password)
	require.NoError(t, err)

	// Correct password should verify
	valid, err := VerifyPassword(password, hash)
	require.NoError(t, err)
	assert.True(t, valid)

	// Wrong password should not verify
	valid, err = VerifyPassword(wrongPassword, hash)
	require.NoError(t, err)
	assert.False(t, valid)
}

func TestVerifyPassword_InvalidHash(t *testing.T) {
	password := "testpassword123"
	invalidHash := "invalid-hash"

	valid, err := VerifyPassword(password, invalidHash)
	assert.Error(t, err)
	assert.False(t, valid)
	assert.Equal(t, ErrInvalidHash, err)
}

func TestHashPasswordWithParams(t *testing.T) {
	password := "testpassword123"

	// Test with custom parameters
	params := &Argon2Params{
		Memory:      32 * 1024, // 32 MB
		Iterations:  2,
		Parallelism: 1,
		SaltLength:  16,
		KeyLength:   32,
	}

	hash, err := HashPasswordWithParams(password, params)
	require.NoError(t, err)
	assert.NotEmpty(t, hash)

	// Should still verify correctly
	valid, err := VerifyPassword(password, hash)
	require.NoError(t, err)
	assert.True(t, valid)
}

func TestDecodeHash(t *testing.T) {
	password := "testpassword123"

	hash, err := HashPassword(password)
	require.NoError(t, err)

	params, salt, key, err := decodeHash(hash)
	require.NoError(t, err)

	assert.Equal(t, DefaultArgon2Params.Memory, params.Memory)
	assert.Equal(t, DefaultArgon2Params.Iterations, params.Iterations)
	assert.Equal(t, DefaultArgon2Params.Parallelism, params.Parallelism)
	assert.Equal(t, DefaultArgon2Params.SaltLength, params.SaltLength)
	assert.Equal(t, DefaultArgon2Params.KeyLength, params.KeyLength)
	assert.Len(t, salt, int(params.SaltLength))
	assert.Len(t, key, int(params.KeyLength))
}

func TestPasswordSecurity(t *testing.T) {
	password := "testpassword123"

	// Same password should generate different hashes
	hash1, err := HashPassword(password)
	require.NoError(t, err)

	hash2, err := HashPassword(password)
	require.NoError(t, err)

	assert.NotEqual(t, hash1, hash2, "Same password should generate different hashes due to random salt")

	// Both should verify correctly
	valid1, err := VerifyPassword(password, hash1)
	require.NoError(t, err)
	assert.True(t, valid1)

	valid2, err := VerifyPassword(password, hash2)
	require.NoError(t, err)
	assert.True(t, valid2)
}

func TestEmptyPassword(t *testing.T) {
	// Empty password should still work (though not recommended)
	hash, err := HashPassword("")
	require.NoError(t, err)

	valid, err := VerifyPassword("", hash)
	require.NoError(t, err)
	assert.True(t, valid)

	// Non-empty password should not verify
	valid, err = VerifyPassword("test", hash)
	require.NoError(t, err)
	assert.False(t, valid)
}

// Benchmark tests
func BenchmarkHashPassword(b *testing.B) {
	password := "testpassword123"

	for i := 0; i < b.N; i++ {
		_, err := HashPassword(password)
		if err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkVerifyPassword(b *testing.B) {
	password := "testpassword123"
	hash, err := HashPassword(password)
	if err != nil {
		b.Fatal(err)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := VerifyPassword(password, hash)
		if err != nil {
			b.Fatal(err)
		}
	}
}
