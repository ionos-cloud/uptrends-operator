package controller

import (
	"context"

	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/manager"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
)

// NewIngressReconciler ...
func NewIngressReconciler(mgr manager.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		Complete(&ingressReconciler{
			Client: mgr.GetClient(),
			scheme: mgr.GetScheme(),
		})
}

type ingressReconciler struct {
	client.Client
	scheme *runtime.Scheme
}

// Reconcile ...
func (s *ingressReconciler) Reconcile(ctx context.Context, r reconcile.Request) (reconcile.Result, error) {
	log := ctrl.LoggerFrom(ctx)
	log.Info("Reconciling Octopinger Config")

	return reconcile.Result{}, nil
}
