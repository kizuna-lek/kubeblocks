package auth

import (
	"fmt"
	"os"
	"path"

	"github.com/spf13/cobra"
	"k8s.io/cli-runtime/pkg/genericclioptions"
	cmdutil "k8s.io/kubectl/pkg/cmd/util"

	"github.com/apecloud/kubeblocks/internal/cli/util"
)

type cleanOptions struct {
	genericclioptions.IOStreams
}

func newCleanCacheCmd(f cmdutil.Factory, streams genericclioptions.IOStreams) *cobra.Command {
	o := &cleanOptions{
		IOStreams: streams,
	}
	cmd := &cobra.Command{
		Use:   "clean-cache",
		Short: "Clean cached auth credentials",
		Args:  cobra.NoArgs,
		Run: func(cmd *cobra.Command, args []string) {
			cmdutil.CheckErr(o.Complete(f, cmd))
			cmdutil.CheckErr(o.Run(f, cmd))
		},
	}
	cmd.AddCommand(newLoginCmd(f, streams))
	return cmd
}

func (o *cleanOptions) Complete(f cmdutil.Factory, cmd *cobra.Command) error {
	return nil
}

func (o *cleanOptions) Run(f cmdutil.Factory, cmd *cobra.Command) error {
	fileDir, err := util.GetCliHomeDir()
	if err != nil {
		return err
	}
	if err := os.Remove(path.Join(fileDir, "credentials")); err != nil {
		return err
	}
	fmt.Fprintln(o.IOStreams.Out, "cache cleaned")
	return nil
}
