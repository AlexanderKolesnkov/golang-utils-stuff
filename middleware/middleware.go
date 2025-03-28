package middleware

import "context"

type AuthCredentials struct {
	authHeader string
}

func NewAuthCredentials(authHeader string) *AuthCredentials {
	return &AuthCredentials{authHeader: authHeader}
}

func (a *AuthCredentials) GetRequestMetadata(ctx context.Context, uri ...string) (map[string]string, error) {
	return map[string]string{
		"authorization": a.authHeader,
	}, nil
}

func (a *AuthCredentials) RequireTransportSecurity() bool {
	return true
}
