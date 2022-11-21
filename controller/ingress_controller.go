package controller

import (
	"context"
	"reflect"
	"strconv"
	"strings"

	"github.com/ionos-cloud/uptrends-operator/api/v1alpha1"
	"github.com/ionos-cloud/uptrends-operator/pkg/utils"

	networkingv1 "k8s.io/api/networking/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/handler"
	"sigs.k8s.io/controller-runtime/pkg/manager"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
	"sigs.k8s.io/controller-runtime/pkg/source"
)

// NewIngressController ...
func NewIngressController(mgr manager.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&networkingv1.Ingress{}).
		Watches(&source.Kind{Type: &networkingv1.Ingress{}}, &handler.EnqueueRequestForObject{}).
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

	items := make(map[string]string)

	for k, v := range in.Annotations {
		if strings.HasPrefix(k, v1alpha1.AnnotationPrefix) {
			items[strings.TrimPrefix(k, v1alpha1.AnnotationPrefix)] = v
		}
	}

	for _, r := range in.Spec.Rules {
		if r.Host == "" {
			continue
		}

		if strings.HasPrefix(r.Host, "*") {
			continue
		}

		monitor := &v1alpha1.Uptrends{
			ObjectMeta: metav1.ObjectMeta{
				Namespace: in.Namespace,
				Name:      r.Host,
			},
			Spec: v1alpha1.UptrendsSpec{
				Type: "HTTPS",
			},
		}

		if v, ok := items["type"]; ok {
			monitor.Spec.Type = v
		}

		monitor.Spec.Interval = 5
		if v, ok := items["interval"]; ok {
			if i, err := strconv.Atoi(v); err == nil {
				monitor.Spec.Interval = i
			}
		}

		if monitor.Spec.Type == "HTTPS" {
			monitor.Spec.Url = "https://" + r.Host
		}

		if monitor.Spec.Type == "HTTP" {
			monitor.Spec.Url = "http://" + r.Host
		}

		existingMonitor := &v1alpha1.Uptrends{}
		if utils.IsObjectFound(ctx, c, in.Namespace, r.Host, existingMonitor) {
			// this is not DaemonSet is not owned by Octopinger
			if ownerRef := metav1.GetControllerOf(existingMonitor); ownerRef == nil || ownerRef.Kind != v1alpha1.CRDResourceKind {
				continue
			}

			if !reflect.DeepEqual(existingMonitor, monitor) {
				existingMonitor = monitor
				err := c.Update(ctx, existingMonitor)
				if err != nil {
					return reconcile.Result{}, err
				}
			}

			continue
		}

		err := c.Create(ctx, monitor)
		if err != nil {
			return reconcile.Result{}, err
		}
	}

	return reconcile.Result{}, nil
}
