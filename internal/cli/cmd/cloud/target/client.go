package target

import (
	"encoding/base64"

	authprovider "github.com/apecloud/kubeblocks/internal/cli/cmd/cloud/auth/provider"
	cloudclientset "github.com/apecloud/kubeblocks/pkg/clientset"

	corev1 "k8s.io/client-go/kubernetes/typed/core/v1"
	"k8s.io/client-go/rest"
	clientcmdapi "k8s.io/client-go/tools/clientcmd/api"
)

type Clientset struct {
	CloudV1alpha1 cloudclientset.Interface
	CoreV1        corev1.CoreV1Interface
}

func New() (*Clientset, error) {

	var client Clientset

	var decoded []byte

	decoded, err := base64.StdEncoding.DecodeString(DBaaSCloudCA)
	if err != nil {
		return nil, err
	}

	config := &rest.Config{
		Host: DBaaSCloudServer,
		TLSClientConfig: rest.TLSClientConfig{
			CAData: decoded,
		},
		// WrapTransport: func(rt http.RoundTripper) http.RoundTripper {
		// 	return nil
		// },
		AuthProvider: &clientcmdapi.AuthProviderConfig{
			Name:   authprovider.ApeOIDCPluginName,
			Config: authprovider.DefaultAuthProviderConfig,
		},
	}

	client.CloudV1alpha1, err = cloudclientset.NewForConfig(config)
	if err != nil {
		return nil, err
	}
	client.CoreV1, err = corev1.NewForConfig(config)
	if err != nil {
		return nil, err
	}

	return &client, nil

	// if err != nil {
	// 	return nil, err
	// }
	// return &Client{RESTClient: cloudClient}

	// restClient, err := rest.RESTClientFor(config)

	// if err != nil {
	// 	return nil, err
	// }

	// cacheDir, err := getDefaultCacheDir()
	// if err != nil {
	// 	return nil, err
	// }

	// httpCacheDir := filepath.Join(cacheDir, "http")

	// discoveryCacheDir := computeDiscoverCacheDir(filepath.Join(cacheDir, "discovery"), config.Host)

	// discoveryClient, err := diskcached.NewCachedDiscoveryClientForConfig(config, discoveryCacheDir, httpCacheDir, 6*time.Hour)

	// if err != nil {
	// 	return nil, err
	// }

	// mapper := restmapper.NewDeferredDiscoveryRESTMapper(discoveryClient)

	// client := &Client{
	// 	RESTClient: restClient,
	// 	RESTMapper: mapper,
	// }

	// return client, nil

}

// func getDefaultCacheDir() (string, error) {
// 	homeDir, err := dbctlutil.GetCliHomeDir()
// 	if err != nil {
// 		return "", err
// 	}
// 	return filepath.Join(homeDir, ".kube", "cache"), nil
// }

// var overlyCautiousIllegalFileCharacters = regexp.MustCompile(`[^(\w/.)]`)

// // computeDiscoverCacheDir takes the parentDir and the host and comes up with a "usually non-colliding" name.
// func computeDiscoverCacheDir(parentDir, host string) string {
// 	// strip the optional scheme from host if its there:
// 	schemelessHost := strings.Replace(strings.Replace(host, "https://", "", 1), "http://", "", 1)
// 	// now do a simple collapse of non-AZ09 characters.  Collisions are possible but unlikely.  Even if we do collide the problem is short lived
// 	safeHost := overlyCautiousIllegalFileCharacters.ReplaceAllString(schemelessHost, "_")
// 	return filepath.Join(parentDir, safeHost)
// }
