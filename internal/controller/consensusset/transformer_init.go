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

package consensusset

import (
	workloads "github.com/apecloud/kubeblocks/apis/workloads/v1alpha1"
	"github.com/apecloud/kubeblocks/internal/controller/graph"
	"github.com/apecloud/kubeblocks/internal/controller/model"
)

type initTransformer struct {
	*workloads.ConsensusSet
}

func (t *initTransformer) Transform(dag *graph.DAG) error {
	vertex := &model.ObjectVertex{
		Obj: t.ConsensusSet,
		OriObj: t.ConsensusSet.DeepCopy(),
	}
	dag.AddVertex(vertex)
	return nil
}