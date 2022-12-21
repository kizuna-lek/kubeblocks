package authprovider

import (
	"github.com/99designs/keyring"
	"github.com/auth0/k8s-pixy-auth/auth"
	"github.com/pkg/errors"
)

type tokenProvider interface {
	GetIDToken() (string, error)
}

func newCachingTokenProviderUsingKeyring(issuer, clientID, audience string, withRefreshToken bool, port uint16, k keyring.Keyring) (tokenProvider, error) {
	atProvider, err := auth.NewDefaultAccessTokenProvider(auth.Issuer{
		IssuerEndpoint: issuer,
		ClientID:       clientID,
		Audience:       audience,
	}, withRefreshToken, port)
	if err != nil {
		return nil, errors.Wrap(err, "could not build access token provider")
	}

	return auth.NewCachingTokenProvider(
		NewKeyringCachingProvider(clientID, audience, k),
		atProvider), nil
}

type TokenReader interface {
	GetToken() (string, error)
}

type reader struct {
	KeyRingReader
}

var tokenStore *reader

func (r *reader) GetToken() (string, error) {

	item, err := r.Get("credentials")
	if err != nil {
		return "", err
	}
	return string(item.Data), nil
}

func GetTokenStore() (TokenReader, error) {
	if tokenStore == nil {
		kr, err := keyringSetup()
		if err != nil {
			return nil, err
		}
		tokenStore = &reader{KeyRingReader: kr}
	}
	return tokenStore, nil
}
