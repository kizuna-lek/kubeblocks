/*
Copyright ApeCloud Inc.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package target

import (
	"context"
	"fmt"

	"github.com/spf13/cobra"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/cli-runtime/pkg/genericclioptions"
	cmdutil "k8s.io/kubectl/pkg/cmd/util"
	"k8s.io/kubectl/pkg/util/templates"

	clouderrors "github.com/apecloud/kubeblocks/internal/cli/cmd/cloud/errors"
	printer "github.com/apecloud/kubeblocks/internal/cli/cmd/cloud/printer"
)

var listExample = templates.Examples(`
		# list all targets
		kbcli cloud target list

		# list a single target with specified NAME
		kbcli cloud target list my-target

		# list a single target in YAML output format
		kbcli cloud target list my-target -o yaml

		# list a single target in JSON output format
		kbcli cloud target list my-target -o json

		# list a single target in wide output format
		kbcli cloud target list my-target -o wide`)

type listOptions struct {
	// TODO(gonglei)
	Options
}

func newListCmd(f cmdutil.Factory, streams genericclioptions.IOStreams) *cobra.Command {
	use := "list"
	alias := "ls"
	o := &listOptions{
		Options: Options{
			IOStreams: streams,
		},
	}
	cmd := &cobra.Command{
		Use:     use,
		Short:   "List all target.",
		Example: listExample,
		Aliases: []string{alias},
		// ValidArgsFunction: utilcomp.ResourceNameCompletionFunc(f, util.GVRToString(types.TargetsGVR())),
		Run: func(cmd *cobra.Command, args []string) {
			clouderrors.CheckErr(o.Validate(args))
			clouderrors.CheckErr(o.Complete())
			clouderrors.CheckErr(o.run(args))
		},
	}
	customFlags(cmd)
	return cmd
}

func customFlags(c *cobra.Command) {
	// TODO(gonglei) add output options, e.g., -o json/yaml/wide
}

// If show-instance, show-component or -o wide is set, output corresponding information,
// if these flags are set on the same time, only one is valid, their priority order is
// show-instance, show-component and -o wide.
func (o *listOptions) run(args []string) error {
	// var printer cluster.Printer

	namespace := o.ID.GetGroups()[0]

	client := o.Cloudclientset

	p := printer.NewTargetPrinter(o.IOStreams.Out)

	// printer = cluster.NewClusterPrinter(c.IOStreams.Out)

	return show(client, namespace, args, o.IOStreams, p)
}

func show(client *Clientset, namespace string, names []string,
	streams genericclioptions.IOStreams, p printer.Printer) error {

	// cluster names are specified by command args
	for _, name := range names {
		target, err := client.CloudV1alpha1.CloudV1alpha1().Targets(namespace).Get(context.TODO(), name, metav1.GetOptions{})
		if err != nil {
			return err
		}
		p.AddRow(target)
	}

	if len(names) > 0 {
		p.Print()
		return nil
	}

	// do not specify any cluster name, we will get all clusters
	targets, err := client.CloudV1alpha1.CloudV1alpha1().Targets(namespace).List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		return err
	}

	// no clusters found
	if len(targets.Items) == 0 {
		fmt.Fprintln(streams.ErrOut, "No resources found")
		return nil
	}

	for _, t := range targets.Items {
		p.AddRow(&t)
	}
	p.Print()
	return nil
}
