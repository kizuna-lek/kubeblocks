package builder

import (
	"context"
	"path/filepath"
	"regexp"
	"strings"
	"time"

	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"k8s.io/apimachinery/pkg/api/meta"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/discovery"
	diskcached "k8s.io/client-go/discovery/cached/disk"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/restmapper"
	"k8s.io/client-go/tools/clientcmd"
	clientcmdapi "k8s.io/client-go/tools/clientcmd/api"
	cmdutil "k8s.io/kubectl/pkg/cmd/util"

	authprovider "github.com/apecloud/kubeblocks/internal/cli/cmd/cloud/auth/provider"
	cloudtarget "github.com/apecloud/kubeblocks/internal/cli/cmd/cloud/target"
	"github.com/apecloud/kubeblocks/internal/cli/util"
	"github.com/apecloud/kubeblocks/pkg/clientset/scheme"
)

var (
	defaultDiscoveryBurst = rest.DefaultBurst
	defaultDiscoveryQPS   = rest.DefaultQPS
)

type buildinConfigFlags struct {
	// genericclioptions.ConfigFlags
	Host     string
	CAData   string
	CacheDir string
}

// overlyCautiousIllegalFileCharacters matches characters that *might* not be supported.  Windows is really restrictive, so this is really restrictive
var overlyCautiousIllegalFileCharacters = regexp.MustCompile(`[^(\w/.)]`)

// computeDiscoverCacheDir takes the parentDir and the host and comes up with a "usually non-colliding" name.
func computeDiscoverCacheDir(parentDir, host string) string {
	// strip the optional scheme from host if its there:
	schemelessHost := strings.Replace(strings.Replace(host, "https://", "", 1), "http://", "", 1)
	// now do a simple collapse of non-AZ09 characters.  Collisions are possible but unlikely.  Even if we do collide the problem is short lived
	safeHost := overlyCautiousIllegalFileCharacters.ReplaceAllString(schemelessHost, "_")
	return filepath.Join(parentDir, safeHost)
}

func (f *buildinConfigFlags) ToDiscoveryClient() (discovery.CachedDiscoveryInterface, error) {
	config, err := f.ToRESTConfig()
	if err != nil {
		return nil, err
	}

	config.Burst = defaultDiscoveryBurst
	config.QPS = defaultDiscoveryQPS

	cacheDir := f.CacheDir

	httpCacheDir := filepath.Join(cacheDir, "http")
	discoveryCacheDir := computeDiscoverCacheDir(filepath.Join(cacheDir, "discovery"), config.Host)

	return diskcached.NewCachedDiscoveryClientForConfig(config, discoveryCacheDir, httpCacheDir, 6*time.Hour)
}

func (f *buildinConfigFlags) ToRESTMapper() (meta.RESTMapper, error) {
	discoveryClient, err := f.ToDiscoveryClient()
	if err != nil {
		return nil, err
	}

	mapper := restmapper.NewDeferredDiscoveryRESTMapper(discoveryClient)
	expander := restmapper.NewShortcutExpander(mapper, discoveryClient)
	return expander, nil
}

func (f *buildinConfigFlags) ToRawKubeConfigLoader() clientcmd.ClientConfig {
	return &clientcmd.DefaultClientConfig
}

func (f *buildinConfigFlags) ToRESTConfig() (*rest.Config, error) {
	// decoded, err := base64.StdEncoding.DecodeString(f.CAData)
	// if err != nil {
	// 	return nil, err
	// }

	config := &rest.Config{
		Host: f.Host,
		TLSClientConfig: rest.TLSClientConfig{
			// TODO(gonglei): CA issue need be handled, https server backend in elb
			// CAData:   decoded,
			Insecure: true,
		},
		ContentConfig: rest.ContentConfig{},
		AuthProvider: &clientcmdapi.AuthProviderConfig{
			Name:   authprovider.ApeOIDCPluginName,
			Config: authprovider.DefaultAuthProviderConfig,
		},
	}
	if err := setKubernetesDefaults(config); err != nil {
		return nil, err
	}
	return config, nil
}

// setKubernetesDefaults sets default values on the provided client config for accessing the
// Kubernetes API or returns an error if any of the defaults are impossible or invalid.
// TODO this isn't what we want.  Each clientset should be setting defaults as it sees fit.
func setKubernetesDefaults(config *rest.Config) error {
	// TODO remove this hack.  This is allowing the GetOptions to be serialized.
	config.GroupVersion = &schema.GroupVersion{Group: "", Version: "v1"}

	if config.APIPath == "" {
		config.APIPath = "/api"
	}
	if config.NegotiatedSerializer == nil {
		// This codec factory ensures the resources are not converted. Therefore, resources
		// will not be round-tripped through internal versions. Defaulting does not happen
		// on the client.
		config.NegotiatedSerializer = scheme.Codecs.WithoutConversion()
	}
	return rest.SetKubernetesDefaults(config)
}

func newBuildinConfigFlags(host, ca, cachedir string) *buildinConfigFlags {
	return &buildinConfigFlags{
		Host:     host,
		CAData:   ca,
		CacheDir: cachedir,
	}

}

func WithTarget(fn CustomCompleteFn) CustomCompleteFn {
	return func(cmd *Command) error {
		factory, err := CompleteTarget(cmd.Cmd)
		if err != nil {
			return err
		}
		if factory != nil {
			cmd.Factory = factory
		}
		if fn != nil {
			return fn(cmd)
		}
		return nil
	}
}

func CompleteTarget(cmd *cobra.Command) (cmdutil.Factory, error) {
	var target string
	targetFlag := cmd.Flags().Lookup("target")
	if targetFlag != nil {
		target = targetFlag.Value.String()

	}
	// fmt.Println(target)

	if len(target) != 0 {
		cloudoption := cloudtarget.Options{}
		if err := cloudoption.Complete(); err != nil {
			return nil, err
		}

		targetclientset := cloudoption.Cloudclientset.CloudV1alpha1.CloudV1alpha1().Targets(cloudoption.OrgName)
		instance, err := targetclientset.Get(context.TODO(), target, metav1.GetOptions{})
		if err != nil {
			return nil, errors.Wrap(err, "error get target")
		}

		// TODO(gonglei): handle target status first, instance data should be used safely
		apiserver := instance.Spec.ServerEndpoints[0].ServerAddress
		ca := instance.Spec.MasterAuth.ClusterCACertificate
		cachedir, err := util.GetCliHomeDir()
		if err != nil {
			return nil, err
		}
		targetclusterfactory := newBuildinConfigFlags(apiserver, ca, cachedir)
		return cmdutil.NewFactory(targetclusterfactory), nil
	}

	return nil, nil

}
