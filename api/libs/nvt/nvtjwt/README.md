# nvtjwt

Sign and verify HS256 JSON Web Tokens with generic claims, on top of [golang-jwt/jwt v5](https://github.com/golang-jwt/jwt).

```go
type Claims struct {
    UserID int64 `json:"uid"`
    jwt.RegisteredClaims
}

token, err := nvtjwt.Encode(Claims{
    UserID: 42,
    RegisteredClaims: jwt.RegisteredClaims{
        ExpiresAt: jwt.NewNumericDate(time.Now().Add(15 * time.Minute)),
    },
}, secret)

claims, err := nvtjwt.Decode(token, secret, &Claims{})
if err != nil {
    // invalid, expired, malformed, wrong algorithm, ...
}
```

| Function | Purpose |
|---|---|
| `Encode(claims, secret string)` | sign claims, return the token |
| `Decode(token, secret string, claims)` | verify the token and fill `claims` |
| `EncodeBytes` / `DecodeBytes` | the same with the secret as `[]byte` |

`Decode` needs a non-nil pointer to your claims struct (`&Claims{}`), or `make(jwt.MapClaims)` for map claims.

## Errors

Match them with `errors.Is`: `jwt.ErrTokenExpired`, `jwt.ErrTokenMalformed`, `jwt.ErrTokenSignatureInvalid`, `jwt.ErrTokenNotValidYet`, and `nvtjwt.ErrEmptySecret` for an empty secret.

## Security

- Only HS256 is accepted; other algorithms, including `none`, are rejected.
- On any error `Decode` returns the zero value of the claims type, never the contents of an unverified token.
- Expiry is checked only when the token has an `exp` claim, so always set `ExpiresAt` when issuing.
- Use a random secret of at least 32 bytes.

## Changed from the previous version

- `Decode` used to return the (possibly forged) claims together with the error; it now returns the zero value.
- The empty-secret error is the exported `ErrEmptySecret`.
- The redundant algorithm check and the unreachable `token.Valid` check were removed; `WithValidMethods` already enforces the algorithm.
