// Package nvthash hashes and verifies passwords with Argon2id, bcrypt, PBKDF2
// and scrypt.
//
// Every algorithm has the same pair of functions. Hash returns a
// self-describing string to store in the database; Check verifies a password
// against it:
//
//	hash, err := nvthash.Argon2idHash(password)
//	if err != nil {
//		return err
//	}
//	// store hash
//
//	ok, err := nvthash.CheckArgon2idHash(password, hash)
//	if err != nil {
//		// the stored hash is malformed (errors.Is(err, nvthash.ErrInvalidHash))
//	}
//	if !ok {
//		// wrong password
//	}
//
// A wrong password is not an error: Check returns false, nil. An error means
// the stored hash cannot be verified.
//
// # Algorithms
//
//	Argon2idHash / CheckArgon2idHash   1 pass, 64 MiB, 4 lanes; the default choice
//	BcryptHash   / CheckBcryptHash     cost 10; passwords are limited to 72 bytes
//	PBKDF2Hash   / CheckPBKDF2Hash     HMAC-SHA256, 600,000 iterations
//	ScryptHash   / CheckScryptHash     N=32768, r=8, p=1
//
// The parameters are written into the hash string, so hashes created with
// other parameters still verify. Salts are random and 16 bytes long, and
// hashes are compared in constant time.
//
// # Layout
//
//	hash.go       ErrInvalidHash and the helpers shared by the PHC-style formats
//	argon2id.go   Argon2id
//	bcrypt.go     bcrypt
//	pbkdf2.go     PBKDF2
//	scrypt.go     scrypt
package nvthash
