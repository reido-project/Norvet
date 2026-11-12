package nvthash

import (
	"errors"
	"fmt"

	"golang.org/x/crypto/bcrypt"
)

const bcryptCost = bcrypt.DefaultCost

// BcryptHash hashes password with bcrypt at the default cost (10) and returns
// the standard "$2a$..." string.
//
// bcrypt only uses the first 72 bytes of a password; a longer password is
// rejected with an error rather than silently truncated.
func BcryptHash(password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcryptCost)
	if err != nil {
		return "", fmt.Errorf("nvthash: bcrypt: %w", err)
	}
	return string(hash), nil
}

// CheckBcryptHash reports whether password matches hash. It returns false, nil
// for a wrong password and an error wrapping ErrInvalidHash for a malformed
// hash.
func CheckBcryptHash(password, hash string) (bool, error) {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	switch {
	case err == nil:
		return true, nil
	case errors.Is(err, bcrypt.ErrMismatchedHashAndPassword):
		return false, nil
	}
	return false, fmt.Errorf("%w: %w", ErrInvalidHash, err)
}
