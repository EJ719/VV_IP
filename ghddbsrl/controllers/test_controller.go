/*
Copyright 2021.

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

package controllers

import (
	"context"

	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"

	testv1 "github.com/ghddbsrl/test-controller/api/v1"
)

// TestReconciler reconciles a Test object
type TestReconciler struct {
	client.Client
	Log logr.Logger
	Scheme *runtime.Scheme
}

//+kubebuilder:rbac:groups=test.my.domain,resources=tests,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=test.my.domain,resources=tests/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=test.my.domain,resources=tests/finalizers,verbs=update

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// TODO(user): Modify the Reconcile function to compare the state specified by
// the Test object against the actual cluster state, and then
// perform operations to make the cluster state reflect the state specified by
// the user.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.8.3/pkg/reconcile
func (r *TestReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	_ = log.FromContext(ctx)

	// your logic here
	
	// Fetch the test instance
	test := &testv1.Test{}
	err := r.Client.Get(context.TODO(), req.NamespacedName, memcached)

	// Check if the Deployment already exists, if not create a new one
	deployment := appsv1.Deployment{}
	err = r.Client.Get(context.TODO(), types.NamespacedName{Name: memcached.Name, Namespace: memcached.Namespace}, deployment)
	if err != nil && errors.IsNotFound(err){
		// Define a new Deployment
		dep := r.deploymentForTest(test)
		reqLogger.Info("Creating a new Deployment.", "Deployment.Namespace", dep.Namespace, "Deployment.Name", dep.Name)
		err = r.Client.Create(context.TODO(), dep)
		if err != nil {
			reqLogger.Error(err, "Failed to create new Deployment.", "Deployment.Namespace", dep.Namespace, "Deployment.Name", dep.Name)
			return ctrl.Result{}, err
		}
		
	}
		return ctrl.Result{}, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *TestReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&testv1.Test{}).
		Complete(r)
}
