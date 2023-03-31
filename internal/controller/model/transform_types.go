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

package model

import (
	"fmt"
	"time"

	snapshotv1 "github.com/kubernetes-csi/external-snapshotter/client/v6/apis/volumesnapshot/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	utilruntime "k8s.io/apimachinery/pkg/util/runtime"
	clientgoscheme "k8s.io/client-go/kubernetes/scheme"
	"sigs.k8s.io/controller-runtime/pkg/client"

	dataprotectionv1alpha1 "github.com/apecloud/kubeblocks/apis/dataprotection/v1alpha1"
	extensionsv1alpha1 "github.com/apecloud/kubeblocks/apis/extensions/v1alpha1"
	workloads "github.com/apecloud/kubeblocks/apis/workloads/v1alpha1"
)

// TODO: copy from lifecycle.transform_types, should replace lifecycle's def

var (
	Scheme = runtime.NewScheme()
)

func init() {
	utilruntime.Must(clientgoscheme.AddToScheme(Scheme))

	utilruntime.Must(workloads.AddToScheme(Scheme))
	utilruntime.Must(dataprotectionv1alpha1.AddToScheme(Scheme))
	utilruntime.Must(snapshotv1.AddToScheme(Scheme))
	utilruntime.Must(extensionsv1alpha1.AddToScheme(Scheme))
}

type Action string

const (
	CREATE = Action("CREATE")
	UPDATE = Action("UPDATE")
	DELETE = Action("DELETE")
	STATUS = Action("STATUS")
)

// RequeueDuration default reconcile requeue after duration
var RequeueDuration = time.Millisecond * 100

type GVKName struct {
	gvk      schema.GroupVersionKind
	ns, name string
}

// ObjectVertex describes expected object spec and how to reach it
// obj always represents the expected part: new object in Create/Update action and old object in Delete action
// oriObj is set in Update action
// all transformers doing their object manipulation works on obj.spec
// the root vertex(i.e. the cluster vertex) will be treated specially:
// as all its meta, spec and status can be updated in one reconciliation loop
// Update is ignored when immutable=true
// orphan object will be force deleted when action is DELETE
type ObjectVertex struct {
	Obj       client.Object
	OriObj    client.Object
	Immutable bool
	IsOrphan  bool
	Action    *Action
}

func (v ObjectVertex) String() string {
	if v.Action == nil {
		return fmt.Sprintf("{obj:%T, name: %s, immutable: %v, orphan: %v, action: nil}",
			v.Obj, v.Obj.GetName(), v.Immutable, v.IsOrphan)
	}
	return fmt.Sprintf("{obj:%T, name: %s, immutable: %v, orphan: %v, action: %v}",
		v.Obj, v.Obj.GetName(), v.Immutable, v.IsOrphan, *v.Action)
}

type ObjectSnapshot map[GVKName]client.Object

type RequeueError interface {
	RequeueAfter() time.Duration
	Reason() string
}

type realRequeueError struct {
	reason       string
	requeueAfter time.Duration
}

func (r *realRequeueError) Error() string {
	return fmt.Sprintf("requeue after: %v as: %s", r.requeueAfter, r.reason)
}

func (r *realRequeueError) RequeueAfter() time.Duration {
	return r.requeueAfter
}

func (r *realRequeueError) Reason() string {
	return r.reason
}