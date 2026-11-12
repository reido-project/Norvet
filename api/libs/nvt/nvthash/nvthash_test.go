package nvthash

import (
	"encoding/hex"
	"errors"
	"strings"
	"testing"
)

type algorithm struct {
	name  string
	hash  func(string) (string, error)
	check func(password, hash string) (bool, error)
}

var algorithms = []algorithm{
	{"argon2id", Argon2idHash, CheckArgon2idHash},
	{"bcrypt", BcryptHash, CheckBcryptHash},
	{"pbkdf2", PBKDF2Hash, CheckPBKDF2Hash},
	{"scrypt", ScryptHash, CheckScryptHash},
}

func TestRoundTrip(t *testing.T) {
	for _, a := range algorithms {
		t.Run(a.name, func(t *testing.T) {
			h1, err := a.hash("correct horse")
			if err != nil {
				t.Fatal(err)
			}
			h2, _ := a.hash("correct horse")
			if h1 == h2 {
				t.Error("two hashes of the same password are identical: the salt is not random")
			}

			if ok, err := a.check("correct horse", h1); !ok || err != nil {
				t.Errorf("right password: ok=%v err=%v", ok, err)
			}
			if ok, err := a.check("wrong horse", h1); ok || err != nil {
				t.Errorf("wrong password: ok=%v err=%v, want false, nil", ok, err)
			}
		})
	}
}

func TestMalformedHashIsAnError(t *testing.T) {
	for _, a := range algorithms {
		for _, bad := range []string{"", "garbage", "$other$x$y", "$$$$$$"} {
			ok, err := a.check("pw", bad)
			if ok || !errors.Is(err, ErrInvalidHash) {
				t.Errorf("%s(%q): ok=%v err=%v, want false and ErrInvalidHash", a.name, bad, ok, err)
			}
		}
	}
}

// An empty hash field must never verify: it would equal the empty output of a
// zero-length derivation.
func TestEmptyHashNeverMatches(t *testing.T) {
	cases := map[string]func(string, string) (bool, error){
		"$argon2id$v=19$m=65536,t=1,p=4$c2FsdHNhbHQ$": CheckArgon2idHash,
		"$pbkdf2-sha256$i=1000$c2FsdHNhbHQ$":          CheckPBKDF2Hash,
		"$scrypt$N=32768,r=8,p=1$c2FsdHNhbHQ$":        CheckScryptHash,
		"$pbkdf2-sha256$i=1000$$aGFzaA":               CheckPBKDF2Hash, // empty salt
		"$argon2id$v=19$m=65536,t=1,p=4$$aGFzaA":      CheckArgon2idHash,
	}
	for encoded, check := range cases {
		ok, err := check("", encoded)
		if ok || !errors.Is(err, ErrInvalidHash) {
			t.Errorf("%q: ok=%v err=%v, want false and ErrInvalidHash", encoded, ok, err)
		}
	}
}

// Zero cost parameters make the underlying libraries panic; they must be
// reported as invalid instead.
func TestZeroCostIsAnErrorNotAPanic(t *testing.T) {
	for _, encoded := range []string{
		"$argon2id$v=19$m=65536,t=0,p=4$c2FsdA$aGFzaA",
		"$argon2id$v=19$m=65536,t=1,p=0$c2FsdA$aGFzaA",
		"$argon2id$v=19$m=65536,t=1,p=256$c2FsdA$aGFzaA",
	} {
		if _, err := CheckArgon2idHash("pw", encoded); !errors.Is(err, ErrInvalidHash) {
			t.Errorf("%q: err=%v, want ErrInvalidHash", encoded, err)
		}
	}
	if _, err := CheckPBKDF2Hash("pw", "$pbkdf2-sha256$i=0$c2FsdA$aGFzaA"); !errors.Is(err, ErrInvalidHash) {
		t.Errorf("pbkdf2 i=0: err=%v, want ErrInvalidHash", err)
	}
}

func TestParamsAreStrict(t *testing.T) {
	for _, bad := range []string{
		"$pbkdf2-sha256$i=10xyz$c2FsdA$aGFzaA",
		"$pbkdf2-sha256$x=10$c2FsdA$aGFzaA",
		"$pbkdf2-sha256$i=-1$c2FsdA$aGFzaA",
		"$pbkdf2-sha256$i=10,i=10$c2FsdA$aGFzaA",
	} {
		if _, err := CheckPBKDF2Hash("pw", bad); !errors.Is(err, ErrInvalidHash) {
			t.Errorf("%q: err=%v, want ErrInvalidHash", bad, err)
		}
	}
}

// PBKDF2-HMAC-SHA256, password "password", salt "salt", 1 iteration, 32 bytes.
func TestPBKDF2KnownVector(t *testing.T) {
	want, _ := hex.DecodeString("120fb6cffcf8b32c43e7225256c4f837a86548c92ccc35480805987cb70be17b")
	encoded := "$pbkdf2-sha256$i=1$" + b64.EncodeToString([]byte("salt")) + "$" + b64.EncodeToString(want)

	if ok, err := CheckPBKDF2Hash("password", encoded); !ok || err != nil {
		t.Errorf("known vector: ok=%v err=%v", ok, err)
	}
	if ok, _ := CheckPBKDF2Hash("Password", encoded); ok {
		t.Error("a different password matched the known vector")
	}
	if !strings.HasPrefix(mustHash(t, PBKDF2Hash), "$pbkdf2-sha256$i=600000$") {
		t.Error("PBKDF2Hash output format changed")
	}
}

func mustHash(t *testing.T, hash func(string) (string, error)) string {
	t.Helper()
	h, err := hash("pw")
	if err != nil {
		t.Fatal(err)
	}
	return h
}
