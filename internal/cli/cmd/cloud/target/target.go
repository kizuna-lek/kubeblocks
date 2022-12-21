package target

import (
	"github.com/spf13/cobra"
	"k8s.io/cli-runtime/pkg/genericclioptions"

	clouderrors "github.com/apecloud/kubeblocks/internal/cli/cmd/cloud/errors"

	cmdutil "k8s.io/kubectl/pkg/cmd/util"
)

type CloudOptions struct {
	genericclioptions.IOStreams
}

func (o *CloudOptions) Validate() error {
	if len(DBaaSCloudCA) == 0 || len(DBaaSCloudServer) == 0 {
		panic("cloud settings invalid")
	}
	return nil
}

func NewTargetCmd(f cmdutil.Factory, streams genericclioptions.IOStreams) *cobra.Command {
	o := &CloudOptions{
		IOStreams: streams,
	}
	cmd := &cobra.Command{
		Use:   "target [command]",
		Short: "targets balabala",
		Long:  "targets balabala",
		Args:  cobra.NoArgs,
		Run: func(cmd *cobra.Command, args []string) {
			cmdutil.CheckErr(o.Validate())
			cmdutil.DefaultSubCommandRun(streams.ErrOut)(cmd, args)
		},
	}
	cmd.AddCommand(
		newRegisterCmd(f, streams),
		newUnRegisterCmd(f, streams),
		newListCmd(f, streams),
	)
	return cmd
}

func newRegisterCmd(f cmdutil.Factory, streams genericclioptions.IOStreams) *cobra.Command {
	o := &registerOption{
		Options: Options{
			IOStreams: streams,
		},
	}
	cmd := &cobra.Command{
		Use:   "register",
		Short: "Register a Kubernetes cluster.",
		Long:  "This command registers a Kubernetes cluster to apecloud.",
		Run: func(cmd *cobra.Command, args []string) {
			clouderrors.CheckErr(o.Validate(args))
			clouderrors.CheckErr(o.Complete(f, cmd, args))
			clouderrors.CheckErr(o.register(f, cmd, args))
		},
	}
	return cmd
}

func newUnRegisterCmd(f cmdutil.Factory, streams genericclioptions.IOStreams) *cobra.Command {
	o := &registerOption{
		Options: Options{
			IOStreams: streams,
		},
	}
	cmd := &cobra.Command{
		Use:   "unregister",
		Short: "Unregister a Kubernetes cluster.",
		Long:  "This command unregister a Kubernetes cluster from apecloud.",
		Run: func(cmd *cobra.Command, args []string) {
			clouderrors.CheckErr(o.Validate(args))
			clouderrors.CheckErr(o.Complete(f, cmd, args))
			clouderrors.CheckErr(o.unregister(f, cmd, args))
		},
	}
	return cmd
}
