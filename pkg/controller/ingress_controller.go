package controller

import (
	"context"
	"fmt"
	"reflect"
	"strconv"
	"strings"

	v1alpha1 "github.com/ionos-cloud/uptrends-operator/api/v1alpha1"
	"github.com/ionos-cloud/uptrends-operator/pkg/finalizers"
	"github.com/ionos-cloud/uptrends-operator/pkg/utils"

	networkingv1 "k8s.io/api/networking/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
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

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
func (c *ingressReconciler) Reconcile(ctx context.Context, r reconcile.Request) (reconcile.Result, error) {
	log := ctrl.LoggerFrom(ctx)
	log.Info("reconcile ingress", "name", r.Name, "namespace", r.Namespace)

	in := &networkingv1.Ingress{}
	err := c.Get(ctx, r.NamespacedName, in)
	if err != nil && errors.IsNotFound(err) {
		// Object not found, return. Created objects are automatically garbage collected.
		log.Info("ingress not found", "name", r.Name, "namespace", r.Namespace)

		return ctrl.Result{}, nil
	}

	if err != nil {
		// Error reading the object - requeue the request.
		return reconcile.Result{}, err
	}

	// Delete if timestamp is set
	if !in.ObjectMeta.DeletionTimestamp.IsZero() {
		if finalizers.HasFinalizer(in, v1alpha1.FinalizerName) {
			return c.reconcileDelete(ctx, in)
		}

		// Delete success
		return ctrl.Result{}, nil
	}

	err = c.reconcileResources(ctx, in)
	if err != nil {
		return reconcile.Result{}, err
	}

	return reconcile.Result{}, nil
}

//nolint:gocyclo
func (c *ingressReconciler) reconcileResources(ctx context.Context, in *networkingv1.Ingress) error {
	existingMonitors := &v1alpha1.UptrendsList{}
	err := c.List(ctx, existingMonitors, client.InNamespace(in.Namespace))
	if err != nil {
		return err
	}

	existingNames := make(map[string]v1alpha1.Uptrends)
	for _, m := range existingMonitors.Items {
		existingNames[m.Name] = m
	}

	annotations := make(map[string]string)

	for k, v := range in.Annotations {
		if strings.HasPrefix(k, v1alpha1.AnnotationPrefix) {
			annotations[strings.TrimPrefix(k, v1alpha1.AnnotationPrefix)] = v
		}
	}

	for _, r := range in.Spec.Rules {
		if r.Host == "" {
			continue
		}

		if strings.HasPrefix(r.Host, "*") {
			continue
		}

		name := fmt.Sprintf("%s-%s", r.Host, in.Name)
		delete(existingNames, name)

		monitor := &v1alpha1.Uptrends{
			ObjectMeta: metav1.ObjectMeta{
				Namespace: in.Namespace,
				Name:      name,
			},
			Spec: v1alpha1.UptrendsSpec{
				Name:        fmt.Sprintf("%s - Uptime", r.Host),
				Interval:    5,
				Type:        "HTTPS",
				Group:       v1alpha1.MonitorGroup{},
				Checkpoints: v1alpha1.MonitorCheckpoints{},
			},
		}

		if v, ok := annotations["type"]; ok {
			monitor.Spec.Type = v
		}

		if v, ok := annotations["interval"]; ok {
			if i, err := strconv.Atoi(v); err == nil {
				monitor.Spec.Interval = i
			}
		}

		if v, ok := annotations["regions"]; ok {
			regions := strings.Split(v, ",")

			for _, r := range regions {
				if i, err := strconv.Atoi(strings.TrimSpace(r)); err == nil {
					monitor.Spec.Checkpoints.Regions = append(monitor.Spec.Checkpoints.Regions, int32(i))
				}
			}
		}

		if v, ok := annotations["checkpoints"]; ok {
			checkpoints := strings.Split(v, ",")
			for _, r := range checkpoints {
				if i, err := strconv.Atoi(strings.TrimSpace(r)); err == nil {
					monitor.Spec.Checkpoints.Checkpoints = append(monitor.Spec.Checkpoints.Checkpoints, int32(i))
				}
			}
		}

		if v, ok := annotations["exclude"]; ok {
			excludes := strings.Split(v, ",")
			for _, r := range excludes {
				if i, err := strconv.Atoi(strings.TrimSpace(r)); err == nil {
					monitor.Spec.Checkpoints.ExcludeCheckpoints = append(monitor.Spec.Checkpoints.ExcludeCheckpoints, int32(i))
				}
			}
		}

		if v, ok := annotations["guid"]; ok {
			monitor.Spec.Group.GUID = v
		}

		if monitor.Spec.Type == "HTTPS" {
			monitor.Spec.Url = "https://" + r.Host
		}

		if monitor.Spec.Type == "HTTP" {
			monitor.Spec.Url = "http://" + r.Host
		}

		existingMonitor := &v1alpha1.Uptrends{}
		if utils.IsObjectFound(ctx, c, in.Namespace, name, existingMonitor) {
			if !reflect.DeepEqual(existingMonitor, monitor) {
				existingMonitor.Spec = monitor.Spec

				err := c.Update(ctx, existingMonitor)
				if err != nil {
					return err
				}
			}

			continue
		}

		err := c.Create(ctx, monitor)
		if err != nil {
			return err
		}
	}

	// clean up
	if len(existingNames) > 0 {
		for _, v := range existingNames {
			err := c.Delete(ctx, &v)
			if err != nil {
				return err
			}
		}
	}

	in.SetFinalizers(finalizers.AddFinalizer(in, v1alpha1.FinalizerName))
	err = c.Update(ctx, in)
	if err != nil && !errors.IsNotFound(err) {
		return err
	}

	return nil
}

//nolint:gocyclo
func (c *ingressReconciler) reconcileDelete(ctx context.Context, in *networkingv1.Ingress) (reconcile.Result, error) {
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

		name := fmt.Sprintf("%s-%s", r.Host, in.Name)

		m := &v1alpha1.Uptrends{}
		err := c.Get(ctx, types.NamespacedName{Namespace: in.Namespace, Name: name}, m)
		if err != nil && errors.IsNotFound(err) {
			continue
		}

		err = c.Delete(ctx, m)
		if err != nil {
			return reconcile.Result{}, err
		}
	}

	in.SetFinalizers(finalizers.RemoveFinalizer(in, v1alpha1.FinalizerName))
	err := c.Update(ctx, in)
	if err != nil && !errors.IsNotFound(err) {
		return reconcile.Result{}, err
	}

	return reconcile.Result{}, nil
}
