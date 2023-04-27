/*
Copyright ApeCloud, Inc.

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

package configmanager

import (
	"context"

	"github.com/fsnotify/fsnotify"

	appsv1alpha1 "github.com/apecloud/kubeblocks/apis/apps/v1alpha1"
)

type ConfigHandler interface {
	OnlineUpdate(ctx context.Context, name string, updatedParams map[string]string) error
	VolumeHandle(ctx context.Context, event fsnotify.Event) error
	MountPoint() []string
}

type ConfigSpecMeta struct {
	*appsv1alpha1.ReloadOptions `json:",inline"`

	ReloadType appsv1alpha1.CfgReloadType       `json:"reloadType"`
	ConfigSpec appsv1alpha1.ComponentConfigSpec `json:"configSpec"`

	ToolConfigs        []appsv1alpha1.ToolConfig
	DownwardAPIOptions []appsv1alpha1.DownwardAPIOption

	// config volume mount path
	TPLConfig  string `json:"tplConfig"`
	MountPoint string `json:"mountPoint"`
	// EngineType string `json:"engineType"`
	// DSN        string `json:"dsn"`

	FormatterConfig appsv1alpha1.FormatterConfig `json:"formatterConfig"`
}

type TPLScriptConfig struct {
	Scripts   string `json:"scripts"`
	FileRegex string `json:"fileRegex"`
	DataType  string `json:"dataType"`
	DSN       string `json:"dsn"`

	FormatterConfig appsv1alpha1.FormatterConfig `json:"formatterConfig"`
}

type ConfigSecondaryRenderMeta struct {
	*appsv1alpha1.ComponentConfigSpec `json:",inline"`

	// secondary template path
	Templates       []string                     `json:"templates"`
	FormatterConfig appsv1alpha1.FormatterConfig `json:"formatterConfig"`
}