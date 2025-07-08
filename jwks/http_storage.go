package jwks

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/MicahParks/jwkset"
)

type Storage interface {
	jwkset.Storage
	Refresh(ctx context.Context) error
}

var _ Storage = (*httpStorage)(nil)

type httpStorage struct {
	*jwkset.MemoryJWKSet
	ctx context.Context
	url string
}

func NewFromHTTP(ctx context.Context, remoteJWKSetURL string, interval time.Duration, errorHandler func(err error)) Storage {
	s := &httpStorage{
		MemoryJWKSet: jwkset.NewMemoryStorage(),
		ctx:          ctx,
		url:          remoteJWKSetURL,
	}

	if interval >= 0 && errorHandler != nil {
		go s.worker(interval, errorHandler)
	}

	return s
}

func (s *httpStorage) Refresh(ctx context.Context) (err error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, s.url, nil)

	if err != nil {
		return fmt.Errorf("failed to create HTTP request for JWK Set refresh: %w", err)
	}

	resp, err := http.DefaultClient.Do(req)

	if err != nil {
		return fmt.Errorf("failed to perform HTTP request for JWK Set refresh: %w", err)
	}

	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("%w: %d", jwkset.ErrInvalidHTTPStatusCode, resp.StatusCode)
	}

	var jwks jwkset.JWKSMarshal

	if err = json.NewDecoder(resp.Body).Decode(&jwks); err != nil {
		return fmt.Errorf("failed to decode JWK Set response: %w", err)
	}

	newSet := make([]jwkset.JWK, len(jwks.Keys))

	for i, marshal := range jwks.Keys {
		jwk, err := jwkset.NewJWKFromMarshal(marshal, jwkset.JWKMarshalOptions{Private: true}, jwkset.JWKValidateOptions{})

		if err != nil {
			return fmt.Errorf("failed to create JWK from JWK Marshal: %w", err)
		}

		newSet[i] = jwk
	}

	if err = s.MemoryJWKSet.KeyReplaceAll(ctx, newSet); err != nil {
		return fmt.Errorf("failed to delete all keys from storage: %w", err)
	}

	return
}

func (s *httpStorage) worker(interval time.Duration, errorHandler func(err error)) {
	tick := time.NewTicker(interval)

	for {
		select {
		case <-s.ctx.Done():
			return
		case <-tick.C:
		}

		if err := s.Refresh(s.ctx); err != nil {
			errorHandler(err)
		}
	}
}
