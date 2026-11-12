package nvthash

import (
	"crypto/rand"
	"encoding/base64"
	"errors"
	"fmt"
	"strconv"
	"strings"
)

// ErrInvalidHash is returned, wrapped with detail, by the Check functions when
// the stored hash is malformed or uses parameters that cannot be verified.
// A password that simply does not match is not an error: Check returns
// false, nil.
var ErrInvalidHash = errors.New("nvthash: invalid hash")

// b64 is the encoding of salts and hashes inside the PHC-style strings.
var b64 = base64.RawStdEncoding

// newSalt returns n random bytes.
func newSalt(n int) ([]byte, error) {
	salt := make([]byte, n)
	if _, err := rand.Read(salt); err != nil {
		return nil, fmt.Errorf("nvthash: generate salt: %w", err)
	}
	return salt, nil
}

// splitHash splits "$id$...$..." into its "$"-separated fields (the first one
// is empty) and checks the algorithm id and the number of fields.
func splitHash(encoded, id string, fields int) ([]string, error) {
	parts := strings.Split(encoded, "$")
	if len(parts) != fields || parts[0] != "" || parts[1] != id {
		return nil, fmt.Errorf("%w: not a %s hash", ErrInvalidHash, id)
	}
	return parts, nil
}

// parseParams parses "k1=v1,k2=v2" into the values, requiring exactly the
// given keys in the given order and plain unsigned decimal numbers.
func parseParams(s string, keys ...string) ([]uint32, error) {
	fields := strings.Split(s, ",")
	if len(fields) != len(keys) {
		return nil, fmt.Errorf("%w: bad parameters %q", ErrInvalidHash, s)
	}

	values := make([]uint32, len(keys))
	for i, field := range fields {
		key, text, ok := strings.Cut(field, "=")
		if !ok || key != keys[i] {
			return nil, fmt.Errorf("%w: bad parameters %q", ErrInvalidHash, s)
		}
		n, err := strconv.ParseUint(text, 10, 32)
		if err != nil {
			return nil, fmt.Errorf("%w: bad parameter %q", ErrInvalidHash, field)
		}
		values[i] = uint32(n)
	}
	return values, nil
}

// decodeSaltAndHash decodes the last two fields of a hash. Both must be
// present: an empty hash would otherwise compare equal to the empty result of
// a zero-length derivation and accept any password.
func decodeSaltAndHash(saltField, hashField string) (salt, hash []byte, err error) {
	salt, err = b64.DecodeString(saltField)
	if err != nil || len(salt) == 0 {
		return nil, nil, fmt.Errorf("%w: bad salt", ErrInvalidHash)
	}
	hash, err = b64.DecodeString(hashField)
	if err != nil || len(hash) == 0 {
		return nil, nil, fmt.Errorf("%w: bad hash", ErrInvalidHash)
	}
	return salt, hash, nil
}
