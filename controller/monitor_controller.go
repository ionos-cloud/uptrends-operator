package controller

import (
	"context"

	"github.com/ionos-cloud/uptrends-operator/api/v1alpha1"
	"sigs.k8s.io/controller-runtime/pkg/handler"
	"sigs.k8s.io/controller-runtime/pkg/manager"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
	"sigs.k8s.io/controller-runtime/pkg/source"

	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

// NewMonitorController ...
func NewMonitorController(mgr manager.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&v1alpha1.Uptrends{}).
		Watches(source.NewKindWithCache(&v1alpha1.Uptrends{}, mgr.GetCache()), &handler.EnqueueRequestForObject{}).
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
	log := ctrl.LoggerFrom(ctx)

	log.Info("reconcile monitor", "name", r.Name, "namespace", r.Namespace)

	return reconcile.Result{}, nil
}
