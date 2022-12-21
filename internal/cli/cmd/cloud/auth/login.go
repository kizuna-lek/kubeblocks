package auth

import (
	"fmt"

	"github.com/spf13/cobra"
	"k8s.io/cli-runtime/pkg/genericclioptions"
	cmdutil "k8s.io/kubectl/pkg/cmd/util"

	authprovider "github.com/apecloud/kubeblocks/internal/cli/cmd/cloud/auth/provider"
)

type LoginOptions struct {
	genericclioptions.IOStreams
}

// --force
// Re-run the web authorization flow even if the given account has valid
// credentials.

func newLoginCmd(f cmdutil.Factory, streams genericclioptions.IOStreams) *cobra.Command {
	o := &LoginOptions{
		IOStreams: streams,
	}
	cmd := &cobra.Command{
		Use:   "login",
		Short: "Authorize dbctl cloud to access the Cloud Platform with ApeCloud user credentials.",
		Long: `Obtains access credentials for your user account via a web-based
authorization flow. When this command completes successfully, it sets the
active account in the current configuration to the account specified. If no
configuration exists, it creates a configuration named default.`,
		Args: cobra.NoArgs,
		Run: func(cmd *cobra.Command, args []string) {
			cmdutil.CheckErr(o.Complete(f, cmd))
			cmdutil.CheckErr(o.Run(f, cmd))
		},
	}
	return cmd

}

func (o *LoginOptions) Complete(f cmdutil.Factory, cmd *cobra.Command) error {
	return nil
}

func (o *LoginOptions) Run(f cmdutil.Factory, cmd *cobra.Command) error {

	if err := Login(); err != nil {
		return err
	}

	fmt.Fprintf(o.IOStreams.Out, "Login success")

	return nil
}

func Login() error {
	if provider, err := authprovider.NewOIDCAuthProvider(); err != nil {
		return err
	} else {
		return provider.Login()
	}
}
