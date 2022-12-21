package auth

import (
	"github.com/spf13/cobra"
	"k8s.io/cli-runtime/pkg/genericclioptions"
	cmdutil "k8s.io/kubectl/pkg/cmd/util"
)

func NewAuthCmd(f cmdutil.Factory, streams genericclioptions.IOStreams) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "auth",
		Short: "Manage oauth2 credentials for the DBaaS Cloud CLI.",
		Long: `The dbctl cloud auth command group lets you grant and revoke authorization to
DBaaS Cloud CLI (dbctl cloud CLI) to access DBaaS Cloud. Typically, when
scripting DBaaS Cloud CLI tools for use on multiple machines, using dbctl cloud
auth activate-service-account is recommended.`,
		Args: cobra.NoArgs,
		Run:  cmdutil.DefaultSubCommandRun(streams.ErrOut),
	}
	cmd.AddCommand(
		newLoginCmd(f, streams),
		newGetCallerIdentityCmd(f, streams),
		newPrintIdentityTokenCmd(f, streams),
		newCleanCacheCmd(f, streams),
	)
	return cmd
}
