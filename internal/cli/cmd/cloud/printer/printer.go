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
	"io"
	"time"

	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/util/duration"

	"github.com/apecloud/kubeblocks/apis/cloud/v1alpha1"
	"github.com/apecloud/kubeblocks/internal/cli/printer"
)

type Printer interface {
	AddRow(obj runtime.Object)
	Print()
}

// TargetPrinter prints cluster info
type TargetPrinter struct {
	tbl *printer.TablePrinter
}

var _ Printer = &TargetPrinter{}

func NewTargetPrinter(out io.Writer) *TargetPrinter {
	p := &TargetPrinter{tbl: printer.NewTablePrinter(out)}
	p.tbl.SetHeader("NAME", "ORGANIZATION", "ENDPOINT", "AGE")
	return p
}

func (p *TargetPrinter) AddRow(obj runtime.Object) {

	target := obj.(*v1alpha1.Target)

	var endpoint string
	if len(target.Spec.ServerEndpoints) > 0 {
		endpoint = target.Spec.ServerEndpoints[0].ServerAddress
	}
	age := duration.HumanDuration(time.Since(target.CreationTimestamp.Time))
	p.tbl.AddRow(target.Name, target.Namespace, endpoint, age)
}

func (p *TargetPrinter) Print() {
	p.tbl.Print()
}
