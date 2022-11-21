package controller

import (
	"context"

	"sigs.k8s.io/controller-runtime/pkg/handler"
	"sigs.k8s.io/controller-runtime/pkg/manager"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
	"sigs.k8s.io/controller-runtime/pkg/source"

	networkingv1 "k8s.io/api/networking/v1"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

// NewMonitorController ...
func NewMonitorController(mgr manager.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&networkingv1.Ingress{}).
		Watches(
			source.NewKindWithCache(&networkingv1.Ingress{}, mgr.GetCache()),
			&handler.EnqueueRequestForObject{}).
		Complete(&ingressReconciler{
			Client: mgr.GetClient(),
			scheme: mgr.GetScheme(),
		})
}

type monitorReconcile struct {
	client.Client
	scheme *runtime.Scheme
}

// Reconcile ...
func (m *monitorReconcile) Reconcile(ctx context.Context, r reconcile.Request) (reconcile.Result, error) {
	return reconcile.Result{}, nil
}
