package controller

import (
	"context"

	"sigs.k8s.io/controller-runtime/pkg/handler"
	"sigs.k8s.io/controller-runtime/pkg/manager"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
	"sigs.k8s.io/controller-runtime/pkg/source"

	networkingv1 "k8s.io/api/networking/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

// NewIngressController ...
func NewIngressController(mgr manager.Manager) error {
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

type ingressReconciler struct {
	client.Client
	scheme *runtime.Scheme
}

// Reconcile ...
func (c *ingressReconciler) Reconcile(ctx context.Context, r reconcile.Request) (reconcile.Result, error) {
	log := ctrl.LoggerFrom(ctx)
	log.Info("reconcile ingress", "name", r.Name, "namespace", r.Namespace)

	in := &networkingv1.Ingress{}
	err := c.Get(ctx, r.NamespacedName, in)
	if err != nil && errors.IsNotFound(err) {
		log.Info("ingress not found", "name", r.Name, "namespace", r.Namespace)

		return ctrl.Result{}, nil
	}

	if err != nil {
		return reconcile.Result{}, err
	}

	return reconcile.Result{}, nil
}
