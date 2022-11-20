package monitor

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	v1alpha1 "github.com/ionos-cloud/uptrends-operator/api/v1alpha1"

	corev1 "k8s.io/api/core/v1"
	networkingv1 "k8s.io/api/networking/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	utils "k8s.io/apimachinery/pkg/util/runtime"
	"k8s.io/apimachinery/pkg/watch"
	"k8s.io/client-go/kubernetes"
	typedcorev1 "k8s.io/client-go/kubernetes/typed/core/v1"
	"k8s.io/client-go/tools/cache"
	"k8s.io/client-go/tools/record"
	"k8s.io/klog/v2"
	"k8s.io/kubectl/pkg/scheme"
)

const (
	controllerName = "uptrends-controller"
)

// ServiceController ...
type ServiceController struct {
	// Kubernetes client internal data structures that watch for events and
	// queue them up. hasSynced is used to synchronize.
	informer  cache.SharedIndexInformer
	hasSynced cache.InformerSynced
	// Cache for all running health checks
	store cache.Store
	// Records kubernetes events
	recorder record.EventRecorder
}

// NewInformer ...
func NewInformer(
	client kubernetes.Interface,
	stop <-chan struct{},
) cache.SharedIndexInformer {
	informer := cache.NewSharedIndexInformer(
		&cache.ListWatch{
			ListFunc: func(options metav1.ListOptions) (runtime.Object, error) {
				return client.NetworkingV1().Ingresses(corev1.NamespaceAll).List(context.TODO(), options)
			},
			WatchFunc: func(options metav1.ListOptions) (watch.Interface, error) {
				return client.NetworkingV1().Ingresses(corev1.NamespaceAll).Watch(context.TODO(), options)
			},
		},
		&networkingv1.Ingress{},
		0, // Skip resync
		cache.Indexers{},
	)

	go informer.Run(stop)

	return informer
}

// NewController ...
func NewController(client kubernetes.Interface, informer cache.SharedIndexInformer) *ServiceController {
	klog.V(4).Info("creating event broadcaster")
	eventBroadcaster := record.NewBroadcaster()
	eventBroadcaster.StartStructuredLogging(0)
	eventBroadcaster.StartRecordingToSink(&typedcorev1.EventSinkImpl{Interface: client.CoreV1().Events("")})
	recorder := eventBroadcaster.NewRecorder(scheme.Scheme, corev1.EventSource{Component: controllerName})

	store := cache.NewStore(cache.MetaNamespaceKeyFunc)

	controller := &ServiceController{
		informer:  informer,
		hasSynced: informer.HasSynced,
		store:     store,
		recorder:  recorder,
	}

	informer.AddEventHandler(cache.ResourceEventHandlerFuncs{
		AddFunc:    controller.enqueue,
		UpdateFunc: controller.handleUpdate,
		DeleteFunc: controller.dequeue,
	})

	return controller
}

// Run ...
func (c *ServiceController) Run(stop <-chan struct{}) error {
	defer utils.HandleCrash()

	klog.Info("starting Ingress controller")

	klog.Info("Waiting for informer caches to sync")
	if ok := cache.WaitForCacheSync(stop, c.hasSynced); !ok {
		return fmt.Errorf("failed to wait for caches to sync")
	}

	klog.Info("started workers")
	<-stop
	klog.Info("shutting down workers")

	return nil
}

func (c *ServiceController) enqueue(obj interface{}) {
	in := obj.(*networkingv1.Ingress)

	klog.Infof("create new monitor: %s", in.Name)

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

		m := v1alpha1.Uptrends{
			ObjectMeta: metav1.ObjectMeta{
				Namespace: in.Namespace,
				Name:      r.Host,
			},
			Spec: v1alpha1.UptrendsSpec{
				Type: "HTTPS",
			},
		}

		if v, ok := items["type"]; ok {
			m.Spec.Type = v
		}

		if v, ok := items["interval"]; ok {
			if i, err := strconv.Atoi(v); err == nil {
				m.Spec.Interval = i
			}
		}

		if m.Spec.Type == "HTTPS" {
			m.Spec.Url = "https://" + r.Host
		}

		if m.Spec.Type == "HTTP" {
			m.Spec.Url = "http://" + r.Host
		}
	}

}

func (c *ServiceController) handleUpdate(old, new interface{}) {

}

func (c *ServiceController) dequeue(obj interface{}) {

}
