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

package resource

import (
	"context"
	"time"

	"github.com/pkg/errors"

	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/manager"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"

	corev1alpha1 "github.com/crossplaneio/crossplane/pkg/apis/core/v1alpha1"
	"github.com/crossplaneio/crossplane/pkg/logging"
)

const (
	controllerNameDefaultClass = "defaultclass.crossplane.io"
	defaultClassWait           = 1 * time.Minute
)

var logDefaultClass = logging.Logger.WithName("controller").WithValues("controller", controllerNameDefaultClass)

// Error strings
const (
	errFailedList             = "unable to list policies for claim kind"
	errFailedPolicyConversion = "unable to convert located policy to coreect kind"
	errNoPolicies             = "unable to locate a policy that specifies a default class for claim kind"
	errMultiplePolicies       = "multiple policies that specify a default class defined for claim kind"
)

// A PolicyKind contains the type metadata for a kind of policy.
type PolicyKind schema.GroupVersionKind

// A PolicyListKind contains the type metadata for a kind of policy.
type PolicyListKind schema.GroupVersionKind

// DefaultClassReconciler reconciles resource claims to the
// default resource class for their given kind according to existing
// policies. Predicates ensure that only claims with no resource class
// reference are reconciled.
type DefaultClassReconciler struct {
	client        client.Client
	newClaim      func() Claim
	newPolicy     func() Policy
	newPolicyList func() PolicyList
}

// NewDefaultClassReconciler creates a new DefaultReconciler for the claim kind
func NewDefaultClassReconciler(m manager.Manager, of ClaimKind, by PolicyKind, byList PolicyListKind) *DefaultClassReconciler {
	nc := func() Claim { return MustCreateObject(schema.GroupVersionKind(of), m.GetScheme()).(Claim) }
	np := func() Policy { return MustCreateObject(schema.GroupVersionKind(by), m.GetScheme()).(Policy) }
	npl := func() PolicyList {
		return MustCreateObject(schema.GroupVersionKind(byList), m.GetScheme()).(PolicyList)
	}

	// Panic early if we've been asked to reconcile a claim, policy, or polict list that has
	// not been registered with our controller manager's scheme.
	_ = nc()
	_ = np()
	_ = npl()

	return &DefaultClassReconciler{
		client:        m.GetClient(),
		newClaim:      nc,
		newPolicy:     np,
		newPolicyList: npl,
	}
}

// Reconcile reconciles a claim to the default class reference for its kind
func (r *DefaultClassReconciler) Reconcile(req reconcile.Request) (reconcile.Result, error) {
	logDefaultClass.V(logging.Debug).Info("Reconciling", "request", req)

	ctx, cancel := context.WithTimeout(context.Background(), reconcileTimeout)
	defer cancel()

	claim := r.newClaim()
	if err := r.client.Get(ctx, req.NamespacedName, claim); err != nil {
		// There's no need to requeue if we no longer exist. Otherwise we'll be
		// requeued implicitly because we return an error.
		return reconcile.Result{}, errors.Wrap(IgnoreNotFound(err), errGetClaim)
	}

	// Get policies for claim kind in claim's namespace
	policies := &unstructured.UnstructuredList{}
	policyList := r.newPolicyList()
	policies.SetGroupVersionKind(policyList.GetObjectKind().GroupVersionKind())
	options := &client.ListOptions{
		Namespace: req.Namespace,
	}
	if err := r.client.List(ctx, options, policies); err != nil {
		// If this is the first time we encounter listing error we'll be
		// requeued implicitly due to the status update. If not, we don't
		// care to requeue because list parameters will not change.
		claim.SetConditions(corev1alpha1.ReconcileError(errors.New(errFailedList)))
		return reconcile.Result{}, errors.Wrap(IgnoreNotFound(r.client.Status().Update(ctx, claim)), errUpdateClaimStatus)
	}

	// Check to see if no defaults defined for claim kind.
	if len(policies.Items) == 0 {
		// If this is the first time we encounter no policies we'll be
		// requeued implicitly due to the status update. If not, we will requeue
		// after a time to see if apolicy has been created.
		claim.SetConditions(corev1alpha1.ReconcileError(errors.New(errNoPolicies)))
		return reconcile.Result{RequeueAfter: defaultClassWait}, errors.Wrap(IgnoreNotFound(r.client.Status().Update(ctx, claim)), errUpdateClaimStatus)
	}

	// Check to see if multiple policies defined for claim kind.
	if len(policies.Items) > 1 {
		// If this is the first time we encounter multiple policies we'll be
		// requeued implicitly due to the status update. If not, we will requeue
		// after a time to see if only one policy class exists.
		claim.SetConditions(corev1alpha1.ReconcileError(errors.New(errMultiplePolicies)))
		return reconcile.Result{RequeueAfter: defaultClassWait}, errors.Wrap(IgnoreNotFound(r.client.Status().Update(ctx, claim)), errUpdateClaimStatus)
	}

	// Make sure single item is of correct policy kind
	policy := r.newPolicy()
	if err := runtime.DefaultUnstructuredConverter.FromUnstructured(policies.Items[0].UnstructuredContent(), policy); err != nil {
		// If this is the first time we encounter conversion error we'll be
		// requeued implicitly due to the status update. If not, we don't
		// care to requeue because conversion will likely not change.
		claim.SetConditions(corev1alpha1.ReconcileError(errors.New(errFailedPolicyConversion)))
		return reconcile.Result{}, errors.Wrap(IgnoreNotFound(r.client.Status().Update(ctx, claim)), errUpdateClaimStatus)
	}

	// Set class reference on claim to default resource class
	claim.SetClassReference(policy.GetDefaultClassReference())

	// Do not requeue, claim controller will see update and claim
	// with class reference set will pass predicates.
	return reconcile.Result{Requeue: false}, errors.Wrap(IgnoreNotFound(r.client.Update(ctx, claim)), errUpdateClaimStatus)
}
