// Package nvtjwt signs and verifies HS256 JSON Web Tokens with a generic,
// type-safe API on top of github.com/golang-jwt/jwt/v5.
//
//	type Claims struct {
//		UserID int64 `json:"uid"`
//		jwt.RegisteredClaims
//	}
//
//	token, err := nvtjwt.Encode(Claims{
//		UserID: 42,
//		RegisteredClaims: jwt.RegisteredClaims{
//			ExpiresAt: jwt.NewNumericDate(time.Now().Add(15 * time.Minute)),
//		},
//	}, secret)
//
//	claims, err := nvtjwt.Decode(token, secret, &Claims{})
//	if err != nil {
//		// invalid, expired, malformed, wrong algorithm, ...
//	}
//	fmt.Println(claims.UserID)
//
// Encode and Decode take the secret as a string; EncodeBytes and DecodeBytes
// take it as bytes.
//
// # Security
//
//   - Only HS256 is accepted. Tokens using any other algorithm, including
//     "none", are rejected.
//   - Decode returns the zero value of the claims type on any error, never the
//     contents of an unverified token.
//   - Expiry is checked only when the token carries an exp claim. Always set
//     ExpiresAt when issuing tokens.
//   - Use a random secret of at least 32 bytes. An empty secret is rejected
//     with ErrEmptySecret, but shorter secrets are not.
//
// # Layout
//
//	jwt.go   Encode, Decode, their byte variants and ErrEmptySecret
package nvtjwt
