package target

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	corev1 "k8s.io/api/core/v1"
	rbacv1 "k8s.io/api/rbac/v1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/util/wait"
	"k8s.io/cli-runtime/pkg/genericclioptions"
	"k8s.io/client-go/dynamic"
	"k8s.io/klog/v2"
	cmdutil "k8s.io/kubectl/pkg/cmd/util"

	"github.com/apecloud/kubeblocks/apis/cloud/v1alpha1"
	cloudauth "github.com/apecloud/kubeblocks/internal/cli/cmd/cloud/auth"
	"github.com/apecloud/kubeblocks/internal/cli/types"
	"github.com/apecloud/kubeblocks/internal/cli/util/helm"
)

// dbctl cloud targets register TARGET_NAME \
//     --context=CONTEXT \ # Target kubeconfig context
//     --kubeconfig=KUBECONFIG # Target kubeconfig

type Options struct {
	targetName     string
	OrgName        string
	ID             cloudauth.Identity
	Cloudclientset *Clientset
	genericclioptions.IOStreams
}

func (o *Options) Complete() error {
	id, err := cloudauth.GetCallerIdentity()
	if err != nil {
		return err
	}
	o.ID = id

	groups := id.GetGroups()
	if len(groups) == 0 {
		return fmt.Errorf("user not belonging to any groups")
	}

	o.OrgName = groups[0]

	o.Cloudclientset, err = New()
	if err != nil {
		return err
	}

	uid := id.GetUID()

	user, err := o.Cloudclientset.CloudV1alpha1.CloudV1alpha1().Users(o.OrgName).Get(context.TODO(), uid, metav1.GetOptions{})

	if err != nil {
		return err
	}

	if !user.Active() {
		return errors.Errorf("account %s is %s", uid, user.Phase())
	}

	return nil
}

func (o *Options) Validate(args []string) error {
	return nil
}

type registerOption struct {
	Options
	kubeclient dynamic.Interface
	installer  *Installer
}

func (o *registerOption) Validate(args []string) error {

	if len(args) == 0 {
		return errors.New("target name is needed")
	}

	o.targetName = args[0]

	return o.Options.Validate(args)
}

func (o *registerOption) Complete(f cmdutil.Factory, cmd *cobra.Command, args []string) error {
	if err := o.Options.Complete(); err != nil {
		return err
	}

	kubeconfig, err := cmd.Flags().GetString("kubeconfig")
	if err != nil {
		return err
	}

	kubecontext, err := cmd.Flags().GetString("context")
	if err != nil {
		return err
	}

	cfg, err := helm.NewActionConfig(types.KubeblockSystemNamespace, kubeconfig, helm.WithContext(kubecontext))
	if err != nil {
		return err
	}

	o.kubeclient, err = f.DynamicClient()
	if err != nil {
		return err
	}

	o.installer = NewInstaller(
		WithClient(o.kubeclient),
		WithNamespace(types.KubeblockSystemNamespace),
		WithHelmConfig(cfg),
		WithVersion(DefaultIdentityVersion),
	)

	return nil
}

func (o *registerOption) register(f cmdutil.Factory, cmd *cobra.Command, args []string) error {

	// 1. get or create target

	namespace := o.OrgName
	uid := o.ID.GetUID()

	user, err := o.Cloudclientset.CloudV1alpha1.CloudV1alpha1().Users(namespace).Get(context.TODO(), uid, metav1.GetOptions{})

	if err != nil {
		return err
	}

	targetclient := o.Cloudclientset.CloudV1alpha1.CloudV1alpha1().Targets(namespace)

	target, err := targetclient.Get(context.TODO(), o.targetName, metav1.GetOptions{})

	if apierrors.IsNotFound(err) {
		// create target if not found
		target = &v1alpha1.Target{}
		target.SetGroupVersionKind(v1alpha1.GroupVersion.WithKind("target"))
		target.SetName(o.targetName)
		target.SetNamespace(namespace)
		now := metav1.Now()
		target.Spec.TokenRefreshTimestamp = &now
		target.OwnerReferences = []metav1.OwnerReference{
			{
				APIVersion: v1alpha1.GroupVersion.String(),
				Kind:       "User",
				Name:       uid,
				UID:        user.GetUID(),
			},
		}
		// TODO(lgong): generate identity service ca bunddle in advance
		target, err = targetclient.Create(context.TODO(), target, metav1.CreateOptions{})
	}

	if err != nil {
		return errors.Wrap(err, "error get target")
	}

	// // 2. wait for secret

	secretclient := o.Cloudclientset.CoreV1.Secrets(namespace)

	waitContext, cancel := context.WithTimeout(context.TODO(), 5*time.Minute)
	defer cancel()

	var token string
	var sec *corev1.Secret

	wait.UntilWithContext(
		waitContext,
		func(ctx context.Context) {
			sec, err = secretclient.Get(waitContext, target.GetName(), metav1.GetOptions{})
			if err != nil {
				return
			}
			token = string(sec.Data["token"])
			cancel()
		},
		200*time.Millisecond,
	)
	if err != nil {
		klog.Errorf("get secret error, err=%s", err.Error())
		return err
	}

	// 3. create indentity service
	// TODO(lgong): using BearerTokenFile for rest.Config to enabling identity-service token rotation
	// TODO(lgong): manage TLS

	o.installer.Sets = []string{
		fmt.Sprintf("apecloud.token=%s", token),
		fmt.Sprintf("apecloud.server=%s", DBaaSCloudServer),
		fmt.Sprintf("apecloud.ca=%s", DBaaSCloudCA),
		fmt.Sprintf("apecloud.targetName=%s", target.GetName()),
		fmt.Sprintf("apecloud.targetNamespace=%s", target.GetNamespace()),
	}

	var notes string

	if notes, err = o.installer.Install(); err != nil {
		return errors.Wrap(err, "failed to install IdentityService")
	}

	// assign roles
	err = createRolebinding(o.kubeclient, user)
	if err != nil {
		return errors.Wrap(err, "failed to assign roles")
	}

	if klog.V(1).Enabled() {
		fmt.Fprintf(o.Out, "IdentityService %s installed\n", o.installer.Version)
		fmt.Fprint(o.Out, notes)
	}
	fmt.Fprintf(o.Out, "register %s finish\n", o.targetName)

	return nil
}

func (o *registerOption) unregister(f cmdutil.Factory, cmd *cobra.Command, args []string) error {

	// uninstall chart
	err := o.installer.Uninstall()
	if err != nil {
		return errors.Wrapf(err, "unregister failed")
	}

	// delete target
	targetclient := o.Cloudclientset.CloudV1alpha1.CloudV1alpha1().Targets(o.OrgName)
	err = targetclient.Delete(context.TODO(), o.targetName, metav1.DeleteOptions{})
	if err != nil && !apierrors.IsNotFound(err) {
		return errors.Wrapf(err, "unregister failed")
	}
	fmt.Fprintf(o.Out, "unregister finish\n")
	return nil
}

func createRolebinding(targetclient dynamic.Interface, user *v1alpha1.User) error {
	var rolebinding rbacv1.ClusterRoleBinding
	kind := rbacv1.SchemeGroupVersion.WithKind("ClusterRoleBinding")
	username := user.GetName()
	rolename := fmt.Sprintf("apecloud-%s", username)
	rolebinding.SetGroupVersionKind(kind)
	rolebinding.SetName(rolename)
	rolebinding.SetNamespace(user.GetNamespace())
	rolebinding.Subjects = []rbacv1.Subject{
		{
			Kind:     rbacv1.UserKind,
			APIGroup: rbacv1.SchemeGroupVersion.Group,
			Name:     fmt.Sprintf("kubeblocks-%s", username),
		},
	}
	role, err := user.GetTargetRoleName()
	if err != nil {
		return err
	}
	rolebinding.RoleRef = rbacv1.RoleRef{
		APIGroup: rbacv1.SchemeGroupVersion.Group,
		Kind:     "ClusterRole",
		Name:     role,
	}

	bt, err := json.Marshal(rolebinding)
	if err != nil {
		return err
	}

	obj := &unstructured.Unstructured{}

	err = json.Unmarshal(bt, obj)
	if err != nil {
		return err
	}

	resource := rbacv1.SchemeGroupVersion.WithResource("clusterrolebindings")
	err = targetclient.Resource(resource).Delete(context.TODO(), rolebinding.GetName(), metav1.DeleteOptions{})
	if err != nil && !apierrors.IsNotFound(err) {
		return err
	}
	_, err = targetclient.Resource(resource).Create(context.TODO(), obj, metav1.CreateOptions{})
	return err
}
