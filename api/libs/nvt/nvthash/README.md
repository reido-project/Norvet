# nvthash

Password hashing and verification with Argon2id, bcrypt, PBKDF2 and scrypt. Every algorithm has the same two functions.

```go
hash, err := nvthash.Argon2idHash(password) // store this string

ok, err := nvthash.CheckArgon2idHash(password, hash)
switch {
case err != nil: // stored hash is malformed: errors.Is(err, nvthash.ErrInvalidHash)
case !ok:        // wrong password
}
```

| Algorithm | Hash | Check | Parameters |
|---|---|---|---|
| Argon2id (default choice) | `Argon2idHash` | `CheckArgon2idHash` | 1 pass, 64 MiB, 4 lanes |
| bcrypt | `BcryptHash` | `CheckBcryptHash` | cost 10, passwords up to 72 bytes |
| PBKDF2 | `PBKDF2Hash` | `CheckPBKDF2Hash` | HMAC-SHA256, 600,000 iterations |
| scrypt | `ScryptHash` | `CheckScryptHash` | N=32768, r=8, p=1 |

All `Check*` functions return `(bool, error)`:

- `true, nil`: the password matches.
- `false, nil`: the password does not match.
- `false, err`: the stored hash is malformed or unverifiable; `errors.Is(err, nvthash.ErrInvalidHash)`.

## Stored format

Argon2id, PBKDF2 and scrypt produce self-describing strings, so hashes made with other parameters still verify:

```
$argon2id$v=19$m=65536,t=1,p=4$<salt>$<hash>
$pbkdf2-sha256$i=600000$<salt>$<hash>
$scrypt$N=32768,r=8,p=1$<salt>$<hash>
```

Salts are random, 16 bytes; salt and hash are unpadded base64. bcrypt uses its standard `$2a$10$...` string.

Verification refuses hashes with an empty salt or hash, or with zero cost parameters, instead of accepting any password or panicking.

## Notes

- bcrypt ignores everything past 72 bytes of the password, so `BcryptHash` returns an error for longer passwords rather than truncating them silently.
- `CheckBcryptHash` used to return only a `bool`; it now returns `(bool, error)` like the others.
