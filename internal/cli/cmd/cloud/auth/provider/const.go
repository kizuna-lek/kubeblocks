package authprovider

var (
	issuerEndpoint string = "https://tenent2.jp.auth0.com/"
	clientID       string = "T8hmxN9SatxiyID4CrgoqYQRRSXq5jaE"
	audience       string = ""
	port           string = "8000"
)

const ApeOIDCPluginName = "ape-oidc"

var DefaultAuthProviderConfig = map[string]string{
	"issuer-url": issuerEndpoint,
	"client-id":  clientID,
	"audience":   audience,
	"port":       port,
}
