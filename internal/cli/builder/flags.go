package builder

import (
	"github.com/spf13/pflag"
	"k8s.io/cli-runtime/pkg/genericclioptions"
)

// GlobalFlags is for flags in global scope
type GlobalFlags struct {
	genericclioptions.RESTClientGetter
	Target string
}

func NewGlobalFlags(delegate genericclioptions.RESTClientGetter) *GlobalFlags {
	return &GlobalFlags{
		RESTClientGetter: delegate,
	}
}

func (f *GlobalFlags) AddFlags(flags *pflag.FlagSet) {
	flags.StringVar(&f.Target, "target", "", "apecloud target to use against, if no target is specified, kubernetes cluster in current kubeconfig context will be used")
}
