package auth

import (
	"context"
	"log"
	"time"

	"github.com/MicahParks/jwkset"
	"github.com/MicahParks/keyfunc/v3"
	"github.com/golang-jwt/jwt/v5"
	"github.com/webmafia/fast"
	"github.com/webmafia/oauth2-authenticator/jwks"
)

type Authenticator struct {
	jwks    jwkset.Storage
	keyfunc keyfunc.Keyfunc
	parser  *jwt.Parser
}

func NewAuthenticator(ctx context.Context, jwksUrl string, refreshInterval time.Duration, issuer string, algs []string) (auth *Authenticator, err error) {
	jwks := jwks.NewFromHTTP(ctx, jwksUrl, refreshInterval, func(err error) {
		log.Println("Error during JWKS refresh:", err)
	})

	// Always load a key into the JWKS before we return.
	if err = jwks.Refresh(ctx); err != nil {
		return
	}

	return NewAuthenticatorWithJWKS(jwks, issuer, algs)
}

func NewAuthenticatorWithJWKS(jwks jwkset.Storage, issuer string, algs []string) (auth *Authenticator, err error) {
	auth = &Authenticator{
		jwks: jwks,
		parser: jwt.NewParser(
			jwt.WithIssuer(issuer),
			jwt.WithValidMethods(algs),
		),
	}

	if auth.keyfunc, err = keyfunc.New(keyfunc.Options{Storage: auth.jwks}); err != nil {
		return
	}

	return
}

func (auth *Authenticator) Validate(token string, dst *Token) (err error) {
	_, err = auth.parser.ParseWithClaims(token, dst, auth.keyfunc.Keyfunc)
	return
}

func (auth *Authenticator) ValidateBytes(token []byte, dst *Token) (err error) {
	return auth.Validate(fast.BytesToString(token), dst)
}

func (auth *Authenticator) ForceRefreshJWKS(ctx context.Context) (err error) {
	if jwks, ok := auth.jwks.(jwks.Storage); ok {
		return jwks.Refresh(ctx)
	}

	return
}
