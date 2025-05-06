# Authenticator
This is a minimal package for parsing and validating tokens.

## Install dependency
`go get github.com/webmafia/oauth2-authenticator`

## Example usage
```go
const url = "https://example.com/.well-known/jwks.json"
const token = "<JWT TOKEN>"
const issuer = "example.com"
var algs = []string{"EdDSA"}

auth, err := NewAuthenticator(context.Background(), url, time.Hour, issuer, algs)

if err != nil {
	return
}

var tok Token

if err = auth.Validate(token, &tok); err != nil {
	return
}

fmt.Printf("Valid token: %#v\n", tok)
```