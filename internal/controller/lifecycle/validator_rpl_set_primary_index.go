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

package lifecycle

import (
	"k8s.io/apimachinery/pkg/types"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"

	appsv1alpha1 "github.com/apecloud/kubeblocks/apis/apps/v1alpha1"
	intctrlutil "github.com/apecloud/kubeblocks/internal/controllerutil"
)

type rplSetPrimaryIndexValidator struct {
	req ctrl.Request
	cli client.Client
	ctx intctrlutil.RequestCtx
}

func (r *rplSetPrimaryIndexValidator) Validate() error {
	cluster := &appsv1alpha1.Cluster{}
	clusterDefinition := &appsv1alpha1.ClusterDefinition{}
	if err := r.cli.Get(r.ctx.Ctx, r.req.NamespacedName, cluster); err != nil {
		return err
	}
	if err := r.cli.Get(r.ctx.Ctx, types.NamespacedName{Name: cluster.Spec.ClusterDefRef}, clusterDefinition); err != nil {
		return err
	}
	// validate primaryIndex and send warning event log necessarily
	if err := cluster.ValidatePrimaryIndex(clusterDefinition); err != nil {
		return newRequeueError(ControllerErrorRequeueTime, err.Error())
	}
	return nil
}