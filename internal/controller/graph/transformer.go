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

package graph

// Transformer transforms a DAG to a new version
type Transformer interface {
	Transform(dag *DAG) error
}

type TransformerChain []Transformer

func (t *TransformerChain) WalkThrough(dag *DAG) error {
	for _, transformer := range *t {
		if err := transformer.Transform(dag); err != nil {
			return err
		}
	}
	return nil
}