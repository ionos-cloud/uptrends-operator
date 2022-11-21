package controller

import (
	"context"

	"github.com/ionos-cloud/uptrends-operator/api/v1alpha1"

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
func NewMonitorController(mgr manager.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&v1alpha1.Uptrends{}).
		Watches(source.NewKindWithCache(&v1alpha1.Uptrends{}, mgr.GetCache()), &handler.EnqueueRequestForObject{}).
		Complete(&monitorReconcile{
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

	err = m.reconcileResource(ctx, mon)
	if err != nil {
		return ctrl.Result{}, err
	}

	return reconcile.Result{}, nil
}

func (m *monitorReconcile) reconcileResource(ctx context.Context, mon *v1alpha1.Uptrends) error {
	auth := context.WithValue(context.Background(), sw.ContextBasicAuth, sw.BasicAuth{
		UserName: "",
		Password: "",
	})

	client := sw.NewAPIClient(sw.NewConfiguration())

	new := sw.Monitor{
		Name:          mon.Spec.Name,
		Url:           mon.Spec.Url,
		MonitorType:   utils.PtrMonitor(sw.MonitorType(mon.Spec.Type)),
		Notes:         mon.Spec.Description,
		CheckInterval: int32(mon.Spec.Interval),
	}

	_, _, err := client.MonitorApi.MonitorPostMonitor(
		auth, &sw.MonitorApiMonitorPostMonitorOpts{
			Monitor: optional.NewInterface(new),
		},
	)
	if err != nil {
		return err
	}

	return nil
}
