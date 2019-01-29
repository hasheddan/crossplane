/*
Copyright 2018 The Crossplane Authors.

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

package provider

import (
	"context"
	"log"

	azurev1alpha1 "github.com/crossplaneio/crossplane/pkg/apis/azure/v1alpha1"
	"github.com/crossplaneio/crossplane/pkg/clients/azure"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/record"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller"
	"sigs.k8s.io/controller-runtime/pkg/handler"
	"sigs.k8s.io/controller-runtime/pkg/manager"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
	"sigs.k8s.io/controller-runtime/pkg/source"
)

const (
	controllerName           = "azure.resourcegroup"
	errorInvalidLocation     = "Invalid location specified"
	errorInvalidSubscription = "Invalid subscription"
)

var (
	ctx           = context.Background()
	result        = reconcile.Result{}
	resultRequeue = reconcile.Result{Requeue: true}
)

// Add creates a new Resource Group Controller and adds it to the Manager with default RBAC. The Manager will set fields on the Controller
// and Start it when the Manager is Started.
func Add(mgr manager.Manager) error {
	return add(mgr, newReconciler(mgr))
}

// Reconciler reconciles a ResourceGroup object
type Reconciler struct {
	client.Client
	scheme     *runtime.Scheme
	kubeclient kubernetes.Interface
	recorder   record.EventRecorder

	// validate func(*google.Credentials, []string) error
	// connect func(*azurev1alpha1.ResourceGroup) (resourcegroup.Client, error)
	create func(*azurev1alpha1.ResourceGroup, azure.Client)
}

// newReconciler returns a new reconcile.Reconciler
func newReconciler(mgr manager.Manager) reconcile.Reconciler {
	r := &Reconciler{
		Client:     mgr.GetClient(),
		scheme:     mgr.GetScheme(),
		kubeclient: kubernetes.NewForConfigOrDie(mgr.GetConfig()),
		recorder:   mgr.GetRecorder(controllerName),
	}
	// r.validate = r._validate
	r.create = r._create
	return r
}

// add adds a new Controller to mgr with r as the reconcile.Reconciler
func add(mgr manager.Manager, r reconcile.Reconciler) error {
	// Create a new controller
	c, err := controller.New(controllerName, mgr, controller.Options{Reconciler: r})
	if err != nil {
		return err
	}

	// Watch for changes to Provider
	err = c.Watch(&source.Kind{Type: &azurev1alpha1.Provider{}}, &handler.EnqueueRequestForObject{})
	if err != nil {
		return err
	}

	return nil
}

// fail - helper function to set fail condition with reason and message
func (r *Reconciler) fail(instance *azurev1alpha1.Provider, reason, msg string) (reconcile.Result, error) {
	instance.Status.UnsetAllConditions()
	instance.Status.SetFailed(reason, msg)
	return resultRequeue, r.Update(context.TODO(), instance)
}

func (r *Reconciler) _create(instance *azurev1alpha1.ResourceGroup, client azure.Client) (reconcile.Result, error) {
	// clusterName := fmt.Sprintf("%s%s", clusterNamePrefix, instance.UID)

	_, err := client.CreateResourceGroup(instance.Spec, client)
	// TODO: better standardizing or errors for Azure
	// if err != nil && !gcp.IsErrorAlreadyExists(err) {
	// 	if gcp.IsErrorBadRequest(err) {
	// 		instance.Status.SetFailed(errorCreateCluster, err.Error())
	// 		// do not requeue on bad requests
	// 		return result, r.Update(ctx, instance)
	// 	}
	// 	return r.fail(instance, errorCreateCluster, err.Error())
	// }
	if err != nil {
		// This is just a placeholder, need better error handling
		return r.fail(instance, errorInvalidSubscription, err.Error())
	}

	// instance.Status.State = gcpcomputev1alpha1.ClusterStateProvisioning

	instance.Status.UnsetAllConditions()
	instance.Status.SetCreating()
	// instance.Status.ClusterName = clusterName

	return resultRequeue, r.Update(ctx, instance)
}

// Reconcile reads that state of the cluster for a ResourceGroup object and makes changes based on the state read
// and what is in the Provider.Spec
// Automatically generate RBAC rules to allow the Controller to read and write Deployments
// +kubebuilder:rbac:groups=apps,resources=deployments,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=azure.crossplane.io,resources=provider,verbs=get;list;watch;create;update;patch;delete
func (r *Reconciler) Reconcile(request reconcile.Request) (reconcile.Result, error) {
	log.Printf("reconciling %s: %v", azurev1alpha1.ProviderKindAPIVersion, request)
	// Fetch the Provider instance
	instance := &azurev1alpha1.ResourceGroup{}
	err := r.Get(ctx, request.NamespacedName, instance)
	if err != nil {
		if errors.IsNotFound(err) {
			// Object not found, return.  Created objects are automatically garbage collected.
			// For additional cleanup logic use finalizers.
			return result, nil
		}
		// Error reading the object - requeue the request.
		return result, err
	}

	// Update status condition
	instance.Status.UnsetAllConditions()
	instance.Status.SetReady()

	return result, r.Update(ctx, instance)
}
