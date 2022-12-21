package authprovider

import (
	"encoding/json"

	"github.com/99designs/keyring"
	"github.com/auth0/k8s-pixy-auth/auth"
	"github.com/pkg/errors"
)

// KeyringCachingProvider satisfies the cachingProvider interface and caches
// tokens using the github.com/99designs/keyring interface
type KeyringCachingProvider struct {
	identifier string
	keyring    auth.KeyringProvider
	// marshalToJSON allows us to mock errors that could happen when
	// marshalling to json
	marshalToJSON func(interface{}) ([]byte, error)
}

func marshalToJSON(toMarshal interface{}) ([]byte, error) {
	marshalled, err := json.Marshal(toMarshal)
	if err != nil {
		return nil, errors.Wrap(err, "could not marshal to json")
	}

	return marshalled, nil
}

// NewKeyringCachingProvider builds a new KeyringCachingProvider using the
// passed in interface satisfiers
func NewKeyringCachingProvider(clientID, audience string, krp auth.KeyringProvider) *KeyringCachingProvider {
	return &KeyringCachingProvider{
		identifier:    "credentials",
		keyring:       krp,
		marshalToJSON: marshalToJSON,
	}
}

// GetTokens gets the TokenResult from keyring
func (kcp *KeyringCachingProvider) GetTokens() (*auth.TokenResult, error) {
	item, err := kcp.keyring.Get(kcp.identifier)
	if err != nil {
		if err == keyring.ErrKeyNotFound {
			return nil, nil
		}

		return nil, errors.Wrap(err, "error getting token information from keyring")
	}

	var tr auth.TokenResult
	err = json.Unmarshal(item.Data, &tr)
	if err != nil {
		return nil, errors.Wrap(err, "could not unmarshal token data from keyring")
	}

	return &tr, nil
}

// CacheTokens stores the TokenResult in keyring
func (kcp *KeyringCachingProvider) CacheTokens(tr *auth.TokenResult) error {
	data, err := kcp.marshalToJSON(tr)
	if err != nil {
		return errors.Wrap(err, "could not marshal token data for caching in keyring")
	}

	err = kcp.keyring.Set(keyring.Item{
		Key:  kcp.identifier,
		Data: data,
	})
	if err != nil {
		return errors.Wrap(err, "error setting token information in keyring")
	}

	return nil
}
