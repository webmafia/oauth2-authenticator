package pkg

import (
	"github.com/gofrs/uuid"
	"github.com/golang-jwt/jwt/v5"
	"github.com/webmafia/fast"
)

var _ jwt.Claims = Token{}

// JWT claims with naming according to: https://www.iana.org/assignments/jwt/jwt.xhtml
type Token struct {
	ID          uuid.UUID       `json:"jti"`
	IssuedAt    jwt.NumericDate `json:"iat"`
	ExpiresAt   jwt.NumericDate `json:"exp"`
	Issuer      string          `json:"iss"`
	Subject     TextualInt      `json:"sub"`
	AuthContext string          `json:"acr"`
	FirstName   string          `json:"given_name"`
	LastName    string          `json:"family_name"`
}

// GetAudience implements jwt.Claims.
func (t Token) GetAudience() (jwt.ClaimStrings, error) {
	return nil, nil
}

// GetExpirationTime implements jwt.Claims.
func (t Token) GetExpirationTime() (*jwt.NumericDate, error) {
	return fast.NoescapeVal(&t.ExpiresAt), nil
}

// GetIssuedAt implements jwt.Claims.
func (t Token) GetIssuedAt() (*jwt.NumericDate, error) {
	return fast.NoescapeVal(&t.IssuedAt), nil
}

// GetIssuer implements jwt.Claims.
func (t Token) GetIssuer() (string, error) {
	return t.Issuer, nil
}

// GetNotBefore implements jwt.Claims.
func (t Token) GetNotBefore() (*jwt.NumericDate, error) {
	return nil, nil
}

// GetSubject implements jwt.Claims.
func (t Token) GetSubject() (string, error) {
	return t.Subject.String(), nil
}
