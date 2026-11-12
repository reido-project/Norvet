package nvthash

import (
	"crypto/subtle"
	"fmt"

	"golang.org/x/crypto/argon2"
)

const (
	argonTime    = 1
	argonMemory  = 64 * 1024 // KiB
	argonThreads = 4
	argonKeyLen  = 32
	argonSaltLen = 16
)

// Argon2idHash hashes password with Argon2id (1 pass, 64 MiB, 4 lanes, 16-byte
// salt, 32-byte key) and returns the PHC string
//
//	$argon2id$v=19$m=65536,t=1,p=4$<salt>$<hash>
//
// Every call uses a fresh random salt.
func Argon2idHash(password string) (string, error) {
	salt, err := newSalt(argonSaltLen)
	if err != nil {
		return "", err
	}

	hash := argon2.IDKey([]byte(password), salt, argonTime, argonMemory, argonThreads, argonKeyLen)

	return fmt.Sprintf("$argon2id$v=%d$m=%d,t=%d,p=%d$%s$%s",
		argon2.Version, argonMemory, argonTime, argonThreads, b64.EncodeToString(salt), b64.EncodeToString(hash)), nil
}

// CheckArgon2idHash reports whether password matches encodedHash, which may
// carry any Argon2id parameters (they are read from the string). It returns
// false, nil for a wrong password and an error wrapping ErrInvalidHash for a
// malformed hash.
func CheckArgon2idHash(password, encodedHash string) (bool, error) {
	parts, err := splitHash(encodedHash, "argon2id", 6)
	if err != nil {
		return false, err
	}

	version, err := parseParams(parts[2], "v")
	if err != nil {
		return false, err
	}
	if version[0] != argon2.Version {
		return false, fmt.Errorf("%w: unsupported argon2 version %d", ErrInvalidHash, version[0])
	}

	params, err := parseParams(parts[3], "m", "t", "p")
	if err != nil {
		return false, err
	}
	memory, passes, threads := params[0], params[1], params[2]
	if passes < 1 || threads < 1 || threads > 255 {
		return false, fmt.Errorf("%w: bad argon2 cost parameters", ErrInvalidHash)
	}

	salt, hash, err := decodeSaltAndHash(parts[4], parts[5])
	if err != nil {
		return false, err
	}

	computed := argon2.IDKey([]byte(password), salt, passes, memory, uint8(threads), uint32(len(hash)))
	return subtle.ConstantTimeCompare(hash, computed) == 1, nil
}
