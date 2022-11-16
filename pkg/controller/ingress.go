package controller

import (
	"context"

	networkingv1 "k8s.io/api/networking/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/tools/cache"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/manager"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
)

// NewIngressReconciler ...
func NewIngressReconciler(mgr manager.Manager) error {
	store := cache.NewStore(cache.MetaNamespaceKeyFunc)

	return ctrl.NewControllerManagedBy(mgr).
		For(&networkingv1.Ingress{}).
		Complete(&ingressReconciler{
			Client: mgr.GetClient(),
			scheme: mgr.GetScheme(),
			store:  store,
		})
}

type ingressReconciler struct {
	client.Client
	scheme *runtime.Scheme
	store  cache.Store
}

// Reconcile ...
func (s *ingressReconciler) Reconcile(ctx context.Context, r reconcile.Request) (reconcile.Result, error) {
	log := ctrl.LoggerFrom(ctx)
	log.Info("reconciling ingress")

	in := &networkingv1.Ingress{}
	err := s.Get(ctx, r.NamespacedName, in)
	if err != nil && errors.IsNotFound(err) {
		// Request object not found, could have been deleted after reconcile request.
		return reconcile.Result{}, nil
	}

	return reconcile.Result{}, nil
}
