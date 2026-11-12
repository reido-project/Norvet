package nvthash

import (
	"crypto/subtle"
	"fmt"

	"golang.org/x/crypto/scrypt"
)

const (
	scryptN       = 32768
	scryptR       = 8
	scryptP       = 1
	scryptKeyLen  = 32
	scryptSaltLen = 16
)

// ScryptHash hashes password with scrypt (N=32768, r=8, p=1, 16-byte salt,
// 32-byte key) and returns the PHC-style string
//
//	$scrypt$N=32768,r=8,p=1$<salt>$<hash>
//
// Every call uses a fresh random salt.
func ScryptHash(password string) (string, error) {
	salt, err := newSalt(scryptSaltLen)
	if err != nil {
		return "", err
	}

	hash, err := scrypt.Key([]byte(password), salt, scryptN, scryptR, scryptP, scryptKeyLen)
	if err != nil {
		return "", fmt.Errorf("nvthash: scrypt: %w", err)
	}

	return fmt.Sprintf("$scrypt$N=%d,r=%d,p=%d$%s$%s", scryptN, scryptR, scryptP, b64.EncodeToString(salt), b64.EncodeToString(hash)), nil
}

// CheckScryptHash reports whether password matches encodedHash, which may
// carry any scrypt parameters (they are read from the string). It returns
// false, nil for a wrong password and an error wrapping ErrInvalidHash for a
// malformed hash.
func CheckScryptHash(password, encodedHash string) (bool, error) {
	parts, err := splitHash(encodedHash, "scrypt", 5)
	if err != nil {
		return false, err
	}

	params, err := parseParams(parts[2], "N", "r", "p")
	if err != nil {
		return false, err
	}

	salt, hash, err := decodeSaltAndHash(parts[3], parts[4])
	if err != nil {
		return false, err
	}

	computed, err := scrypt.Key([]byte(password), salt, int(params[0]), int(params[1]), int(params[2]), len(hash))
	if err != nil {
		return false, fmt.Errorf("%w: %w", ErrInvalidHash, err)
	}
	return subtle.ConstantTimeCompare(hash, computed) == 1, nil
}
