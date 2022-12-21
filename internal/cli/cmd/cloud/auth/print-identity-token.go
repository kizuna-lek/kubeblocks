package auth

import (
	"encoding/json"
	"fmt"

	"github.com/99designs/keyring"
	"github.com/spf13/cobra"
	"k8s.io/cli-runtime/pkg/genericclioptions"
	cmdutil "k8s.io/kubectl/pkg/cmd/util"

	authprovider "github.com/apecloud/kubeblocks/internal/cli/cmd/cloud/auth/provider"
)

type identityPrintOptions struct {
	genericclioptions.IOStreams
}

func newPrintIdentityTokenCmd(f cmdutil.Factory, streams genericclioptions.IOStreams) *cobra.Command {
	o := &identityPrintOptions{
		IOStreams: streams,
	}

	cmd := &cobra.Command{
		Use:   "print-identity-token",
		Short: "Print an identity token for the credentialed the specified account.",
		Args:  cobra.NoArgs,
		Run: func(cmd *cobra.Command, args []string) {
			cmdutil.CheckErr(o.Complete(f, cmd))
			cmdutil.CheckErr(o.Run(f, cmd))
		},
	}
	cmd.AddCommand(newLoginCmd(f, streams))
	return cmd
}

func (o *identityPrintOptions) Complete(f cmdutil.Factory, cmd *cobra.Command) error {
	return nil
}

func (o *identityPrintOptions) Run(f cmdutil.Factory, cmd *cobra.Command) error {

	reader, err := authprovider.GetTokenStore()
	if err != nil {
		return err
	}

	tokenData, err := getTokenFromCache(reader)

	if err != nil {
		return err
	}

	var tk token
	err = json.Unmarshal([]byte(tokenData), &tk)
	if err != nil {
		return err
	}

	fmt.Fprintln(o.IOStreams.Out, tk.IDToken)

	return nil
}

func getTokenFromCache(reader authprovider.TokenReader) (string, error) {

	if reader == nil {
		panic("nil pointer")
	}

	tokenData, err := reader.GetToken()

	if err == nil {
		return tokenData, nil
	}

	if err != keyring.ErrKeyNotFound {
		return tokenData, err
	}
	// token not found in cache, login
	if err = Login(); err != nil {
		return "", err
	}
	return reader.GetToken()
}
