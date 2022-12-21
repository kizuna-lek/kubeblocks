package cloud

import (
	"github.com/spf13/cobra"
	"k8s.io/cli-runtime/pkg/genericclioptions"
	cmdutil "k8s.io/kubectl/pkg/cmd/util"

	"github.com/apecloud/kubeblocks/internal/cli/cmd/cloud/auth"
	"github.com/apecloud/kubeblocks/internal/cli/cmd/cloud/target"
)

func NewCloudCmd(f cmdutil.Factory, streams genericclioptions.IOStreams) *cobra.Command {
	rootCmd := &cobra.Command{
		Use:   "cloud balabala",
		Short: "manage DBaaS Cloud resources and developer workflow",
		Long:  `The DBaaS Cloud CLI manages authentication, local configuration, developer workflow, and interactions with the DBaaS Cloud APIs.`,
		Run:   cmdutil.DefaultSubCommandRun(streams.ErrOut),
	}
	rootCmd.AddCommand(
		target.NewTargetCmd(f, streams),
		auth.NewAuthCmd(f, streams),
	)
	return rootCmd
}
