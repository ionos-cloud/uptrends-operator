package controller

import (
	"context"
	"net/http"

	"github.com/ionos-cloud/uptrends-operator/api/v1alpha1"
	"github.com/ionos-cloud/uptrends-operator/pkg/credentials"
	"github.com/ionos-cloud/uptrends-operator/pkg/finalizers"

	"github.com/antihax/optional"
	sw "github.com/ionos-cloud/uptrends-go"
	"github.com/ionos-cloud/uptrends-go/utils"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/handler"
	"sigs.k8s.io/controller-runtime/pkg/manager"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
	"sigs.k8s.io/controller-runtime/pkg/source"
)

// NewMonitorController ...
func NewMonitorController(mgr manager.Manager, creds *credentials.API) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&v1alpha1.Uptrends{}).
		Watches(source.NewKindWithCache(&v1alpha1.Uptrends{}, mgr.GetCache()), &handler.EnqueueRequestForObject{}).
		Complete(&monitorReconcile{
			Client: mgr.GetClient(),
			scheme: mgr.GetScheme(),
			creds:  creds,
		})
}

type monitorReconcile struct {
	client.Client
	creds  *credentials.API
	scheme *runtime.Scheme
}

// Reconcile ...
func (m *monitorReconcile) Reconcile(ctx context.Context, r reconcile.Request) (reconcile.Result, error) {
	log := ctrl.LoggerFrom(ctx)
	log.Info("reconcile monitor", "name", r.Name, "namespace", r.Namespace)

	mon := &v1alpha1.Uptrends{}
	err := m.Get(ctx, r.NamespacedName, mon)
	if err != nil && errors.IsNotFound(err) {
		// Object not found, return. Created objects are automatically garbage collected.
		log.Info("monitor not found", "name", r.Name, "namespace", r.Namespace)

		return ctrl.Result{}, nil
	}

	if err != nil {
		// Error reading the object - requeue the request.
		return reconcile.Result{}, err
	}

	// Delete if timestamp is set
	if !mon.ObjectMeta.DeletionTimestamp.IsZero() {
		if finalizers.HasFinalizer(mon, v1alpha1.FinalizerName) {
			err := m.reconcileDelete(ctx, mon)
			if err != nil {
				return ctrl.Result{}, err
			}
		}

		// Delete
		return ctrl.Result{}, nil
	}

	err = m.reconcileResources(ctx, mon)
	if err != nil {
		// Error reconciling uptrends sub-resources - requeue the request.
		return reconcile.Result{}, err
	}

	return reconcile.Result{}, nil
}

func (m *monitorReconcile) reconcileResources(ctx context.Context, uptrends *v1alpha1.Uptrends) error {
	err := m.reconcileStatus(ctx, uptrends)
	if err != nil {
		return err
	}

	err = m.reconcileMonitor(ctx, uptrends)
	if err != nil {
		return err
	}

	return nil
}

func (m *monitorReconcile) reconcileDelete(ctx context.Context, mon *v1alpha1.Uptrends) error {
	auth := context.WithValue(ctx, sw.ContextBasicAuth, sw.BasicAuth{
		UserName: m.creds.Username,
		Password: m.creds.Password,
	})

	client := sw.NewAPIClient(sw.NewConfiguration())

	resp, err := client.MonitorApi.MonitorDeleteMonitor(auth, mon.Status.MonitorGuid)
	if err != nil && resp.StatusCode != http.StatusNotFound { // assume that this was already deleted
		return err
	}

	mon.SetFinalizers(finalizers.RemoveFinalizer(mon, v1alpha1.FinalizerName))
	err = m.Update(ctx, mon)
	if err != nil && !errors.IsNotFound(err) {
		return err
	}

	return nil
}

func (m *monitorReconcile) reconcileStatus(ctx context.Context, uptrends *v1alpha1.Uptrends) error {
	phase := v1alpha1.UptrendsPhaseNone

	if uptrends.Status.MonitorGuid != "" {
		phase = v1alpha1.UptrendsPhaseRunning
	}

	if uptrends.Status.Phase != phase {
		uptrends.Status.Phase = phase

		return m.Status().Update(ctx, uptrends)
	}

	return nil
}

func (m *monitorReconcile) reconcileMonitor(ctx context.Context, mon *v1alpha1.Uptrends) error {
	if mon.Status.Phase == v1alpha1.UptrendsPhaseRunning {
		return nil
	}

	auth := context.WithValue(ctx, sw.ContextBasicAuth, sw.BasicAuth{
		UserName: m.creds.Username,
		Password: m.creds.Password,
	})

	client := sw.NewAPIClient(sw.NewConfiguration())

	new := sw.Monitor{
		Name:          mon.Spec.Name,
		Url:           mon.Spec.Url,
		MonitorType:   utils.PtrMonitor(sw.MonitorType(mon.Spec.Type)),
		Notes:         mon.Spec.Description,
		CheckInterval: int32(mon.Spec.Interval),
	}

	up, _, err := client.MonitorApi.MonitorPostMonitor(
		auth, &sw.MonitorApiMonitorPostMonitorOpts{
			Monitor: optional.NewInterface(new),
		},
	)
	if err != nil {
		return err
	}

	mon.SetFinalizers(finalizers.AddFinalizer(mon, v1alpha1.FinalizerName))
	err = m.Update(ctx, mon)
	if err != nil && !errors.IsNotFound(err) {
		return err
	}

	mon.Status.MonitorGuid = up.MonitorGuid
	err = m.Status().Update(ctx, mon)
	if err != nil {
		return err
	}

	return nil
}
