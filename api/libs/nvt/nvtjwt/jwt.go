package nvtjwt

import (
	"errors"

	"github.com/golang-jwt/jwt/v5"
)

// ErrEmptySecret is returned when the signing secret is empty.
var ErrEmptySecret = errors.New("nvtjwt: empty JWT secret")

// signingMethod is the only algorithm this package signs and accepts.
var signingMethod = jwt.SigningMethodHS256

// Encode signs claims with secret (HS256) and returns the compact token.
//
// claims can be any jwt.Claims: a struct embedding jwt.RegisteredClaims, a
// jwt.RegisteredClaims, or a jwt.MapClaims. It returns ErrEmptySecret when
// secret is empty.
func Encode[T jwt.Claims](claims T, secret string) (string, error) {
	return EncodeBytes(claims, []byte(secret))
}

// EncodeBytes is Encode with the secret as bytes.
func EncodeBytes[T jwt.Claims](claims T, secret []byte) (string, error) {
	if len(secret) == 0 {
		return "", ErrEmptySecret
	}

	return jwt.NewWithClaims(signingMethod, claims).SignedString(secret)
}

// Decode verifies tokenString against secret (HS256 only), fills claims from
// it and returns claims.
//
// claims must be a non-nil pointer to a claims struct (for example
// &MyClaims{}) or a non-nil jwt.MapClaims (make(jwt.MapClaims)). The
// registered time claims (exp, nbf, iat) are checked when present; a token
// without exp never expires.
//
// On any error Decode returns the zero value of T, never the claims of an
// unverified token, so a caller that ignores the error cannot act on forged
// data. Errors can be matched with errors.Is against the jwt package's
// sentinels (jwt.ErrTokenExpired, jwt.ErrTokenMalformed,
// jwt.ErrTokenSignatureInvalid, ...), and ErrEmptySecret when secret is empty.
func Decode[T jwt.Claims](tokenString string, secret string, claims T) (T, error) {
	return DecodeBytes(tokenString, []byte(secret), claims)
}

// DecodeBytes is Decode with the secret as bytes.
func DecodeBytes[T jwt.Claims](tokenString string, secret []byte, claims T) (T, error) {
	var zero T

	if len(secret) == 0 {
		return zero, ErrEmptySecret
	}

	_, err := jwt.ParseWithClaims(tokenString, claims, func(*jwt.Token) (any, error) {
		return secret, nil
	}, jwt.WithValidMethods([]string{signingMethod.Alg()}))
	if err != nil {
		return zero, err
	}

	return claims, nil
}
