package authprovider

import (
	"fmt"
	"net/http"

	"github.com/pkg/errors"
	"github.com/spf13/cast"
	"k8s.io/apimachinery/pkg/util/net"
	restclient "k8s.io/client-go/rest"
	"k8s.io/klog/v2"
)

const (
	cfgIssuerURL   = "issuer-url"
	cfgClientID    = "client-id"
	cfgAudience    = "audience"
	cfgPort        = "port"
	cfgDefaultPort = 8000
)

func init() {
	if err := restclient.RegisterAuthProviderPlugin(ApeOIDCPluginName, newOIDCAuthProvider); err != nil {
		klog.Fatalf("Failed to register oidc auth plugin: %v", err)
	}
}

type oidcAuthProvider struct {
	tokenProvider
}

func (p *oidcAuthProvider) WrapTransport(rt http.RoundTripper) http.RoundTripper {
	return &roundTripper{
		wrapped:  rt,
		provider: p,
	}
}

// Login into idp, cache enabled by default
func (p *oidcAuthProvider) Login() error {
	_, err := p.GetIDToken()
	// fmt.Fprintf(os.Stdout, "login success, token=%s", token)
	return errors.Wrap(err, "auth error")
}

func newOIDCAuthProvider(clusterAddress string, config map[string]string, persister restclient.AuthProviderConfigPersister) (restclient.AuthProvider, error) {

	issuer := config[cfgIssuerURL]
	if issuer == "" {
		return nil, fmt.Errorf("must provide %s", cfgIssuerURL)
	}

	clientID := config[cfgClientID]
	if clientID == "" {
		return nil, fmt.Errorf("must provide %s", cfgClientID)
	}

	port := cast.ToUint16(config[cfgPort]) | cfgDefaultPort

	allowRefresh := false

	k, err := keyringSetup()
	if err != nil {
		return nil, errors.Wrap(err, "could not set up keyring")
	}

	provider, err := newCachingTokenProviderUsingKeyring(issuer, clientID, audience, allowRefresh, port, k)
	if err != nil {
		return nil, errors.Wrap(err, "could not build caching token provider")
	}
	return &oidcAuthProvider{provider}, nil

}

func NewOIDCAuthProvider() (restclient.AuthProvider, error) {
	return newOIDCAuthProvider("", DefaultAuthProviderConfig, nil)
}

type roundTripper struct {
	provider *oidcAuthProvider
	wrapped  http.RoundTripper
}

var _ net.RoundTripperWrapper = &roundTripper{}

func (r *roundTripper) RoundTrip(req *http.Request) (*http.Response, error) {
	if len(req.Header.Get("Authorization")) != 0 {
		return r.wrapped.RoundTrip(req)
	}
	token, err := r.provider.GetIDToken()
	if err != nil {
		return nil, errors.Wrap(err, "auth error")
	}
	// klog.Infof("get id token, token=%s", token)

	// shallow copy of the struct
	r2 := new(http.Request)
	*r2 = *req
	// deep copy of the Header so we don't modify the original
	// request's Header (as per RoundTripper contract).
	r2.Header = make(http.Header)
	for k, s := range req.Header {
		r2.Header[k] = s
	}
	r2.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))

	return r.wrapped.RoundTrip(r2)
}

func (r *roundTripper) WrappedRoundTripper() http.RoundTripper { return r.wrapped }
