package nvthash

import (
	"crypto/pbkdf2"
	"crypto/sha256"
	"crypto/subtle"
	"fmt"
)

const (
	pbkdf2Iterations = 600000
	pbkdf2KeyLen     = 32
	pbkdf2SaltLen    = 16
)

// PBKDF2Hash hashes password with PBKDF2-HMAC-SHA256 (600,000 iterations,
// 16-byte salt, 32-byte key) and returns the PHC-style string
//
//	$pbkdf2-sha256$i=600000$<salt>$<hash>
//
// Every call uses a fresh random salt.
func PBKDF2Hash(password string) (string, error) {
	salt, err := newSalt(pbkdf2SaltLen)
	if err != nil {
		return "", err
	}

	hash, err := pbkdf2.Key(sha256.New, password, salt, pbkdf2Iterations, pbkdf2KeyLen)
	if err != nil {
		return "", fmt.Errorf("nvthash: pbkdf2: %w", err)
	}

	return fmt.Sprintf("$pbkdf2-sha256$i=%d$%s$%s", pbkdf2Iterations, b64.EncodeToString(salt), b64.EncodeToString(hash)), nil
}

// CheckPBKDF2Hash reports whether password matches encodedHash, which may
// carry any iteration count (it is read from the string). It returns false,
// nil for a wrong password and an error wrapping ErrInvalidHash for a
// malformed hash.
func CheckPBKDF2Hash(password, encodedHash string) (bool, error) {
	parts, err := splitHash(encodedHash, "pbkdf2-sha256", 5)
	if err != nil {
		return false, err
	}

	params, err := parseParams(parts[2], "i")
	if err != nil {
		return false, err
	}
	if params[0] < 1 {
		return false, fmt.Errorf("%w: bad pbkdf2 iteration count", ErrInvalidHash)
	}

	salt, hash, err := decodeSaltAndHash(parts[3], parts[4])
	if err != nil {
		return false, err
	}

	computed, err := pbkdf2.Key(sha256.New, password, salt, int(params[0]), len(hash))
	if err != nil {
		return false, fmt.Errorf("%w: %w", ErrInvalidHash, err)
	}
	return subtle.ConstantTimeCompare(hash, computed) == 1, nil
}
