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
	"strings"

	"helm.sh/helm/v3/pkg/action"
	"helm.sh/helm/v3/pkg/repo"
	"k8s.io/client-go/dynamic"

	"github.com/apecloud/kubeblocks/internal/cli/types"
	"github.com/apecloud/kubeblocks/internal/cli/util/helm"
)

type Installer struct {
	HelmCfg *action.Configuration

	Namespace string
	Version   string
	Sets      []string
	client    dynamic.Interface
}

type Option func(*Installer)

func (i *Installer) Install() (string, error) {
	// Add repo, if exits, will update it
	if err := helm.AddRepo(&repo.Entry{
		Name: types.ApeCloudChartName,
		URL:  types.ApeCloudChartURL,
	}); err != nil {
		return "", err
	}

	var sets []string
	for _, set := range i.Sets {
		splitSet := strings.Split(set, ",")
		sets = append(sets, splitSet...)
	}
	chart := helm.InstallOpts{
		Name:      types.ApeCloudChartName,
		Chart:     types.ApeCloudChartName + "/" + types.ApeCloudIdentityService,
		Wait:      true,
		Version:   i.Version,
		Namespace: i.Namespace,
		Sets:      sets,
		Login:     true,
		TryTimes:  2,
	}

	notes, err := chart.Install(i.HelmCfg)
	if err != nil {
		return "", err
	}

	return notes, nil
}

// Uninstall remove KubeBlocks
func (i *Installer) Uninstall() error {
	chart := helm.InstallOpts{
		Name:      types.ApeCloudChartName,
		Namespace: i.Namespace,
	}

	if err := chart.UnInstall(i.HelmCfg); err != nil {
		return err
	}

	if err := helm.RemoveRepo(&repo.Entry{Name: types.ApeCloudChartName, URL: types.ApeCloudChartURL}); err != nil {
		return err
	}

	return nil
}

func NewInstaller(opts ...Option) *Installer {
	installer := &Installer{}
	for _, opt := range opts {
		opt(installer)
	}
	return installer
}

func WithHelmConfig(cfg *action.Configuration) Option {
	return func(i *Installer) {
		i.HelmCfg = cfg
	}
}

func WithNamespace(namespace string) Option {
	return func(i *Installer) {
		i.Namespace = namespace
	}
}

func WithVersion(version string) Option {
	return func(i *Installer) {
		i.Version = version
	}
}

func WithSets(sets []string) Option {
	return func(i *Installer) {
		i.Sets = sets
	}
}

func WithClient(client dynamic.Interface) Option {
	return func(i *Installer) {
		i.client = client
	}
}
