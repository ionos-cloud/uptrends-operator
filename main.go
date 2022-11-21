package main

import (
	"context"
	"fmt"

	goruntime "runtime"

	"github.com/ionos-cloud/uptrends-operator/controller"
	"github.com/ionos-cloud/uptrends-operator/pkg/utils"

	"github.com/spf13/cobra"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/tools/leaderelection/resourcelock"
	"k8s.io/klog/v2"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/healthz"
	"sigs.k8s.io/controller-runtime/pkg/log/zap"
	"sigs.k8s.io/controller-runtime/pkg/manager"
	//+kubebuilder:scaffold:imports
)

type flags struct {
	KubeConfig           string
	MasterURL            string
	MetricsAddr          string
	ProbeAddr            string
	EnableLeaderElection bool
}

var f = &flags{}

var (
	scheme   = runtime.NewScheme()
	setupLog = ctrl.Log.WithName("setup")
)

var rootCmd = &cobra.Command{
	Use: "controller",
	RunE: func(cmd *cobra.Command, args []string) error {
		return run(cmd.Context())
	},
}

func printVersion() {
	setupLog.Info(fmt.Sprintf("Go Version: %s", goruntime.Version()))
	setupLog.Info(fmt.Sprintf("Go OS/Arch: %s/%s", goruntime.GOOS, goruntime.GOARCH))
}

func init() {
	rootCmd.Flags().BoolVar(&f.EnableLeaderElection, "leader-elect", f.EnableLeaderElection, "only one controller")
	rootCmd.Flags().StringVar(&f.KubeConfig, "kubeconfig", f.KubeConfig, "Path to a kubeconfig. Only required if out-of-cluster.")
	rootCmd.Flags().StringVar(&f.MasterURL, "master", f.MasterURL, "The address of the Kubernetes API server. Overrides any value in kubeconfig. Only required if out-of-cluster.")
	rootCmd.Flags().StringVar(&f.MetricsAddr, "metrics-bind-address", ":8080", "metrics endpoint")
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		klog.Error(err, "unable to run controller")
	}
}

func run(ctx context.Context) error {
	opts := zap.Options{
		Development: true,
	}

	ctrl.SetLogger(zap.New(zap.UseFlagOptions(&opts)))

	printVersion()

	options := manager.Options{
		BaseContext:                func() context.Context { return ctx },
		HealthProbeBindAddress:     f.ProbeAddr,
		LeaderElection:             f.EnableLeaderElection,
		LeaderElectionID:           "j8yhqdnj.uptrends.ionos-cloud.github.io",
		LeaderElectionResourceLock: resourcelock.LeasesResourceLock,
		MetricsBindAddress:         f.MetricsAddr,
		Namespace:                  "",
		NewClient:                  utils.DefaultNewClientWithMetrics,
		Port:                       9443,
		Scheme:                     scheme,
	}

	mgr, err := ctrl.NewManager(ctrl.GetConfigOrDie(), options)
	if err != nil {
		return err
	}

	err = setupControllers(f, mgr)
	if err != nil {
		return err
	}

	//+kubebuilder:scaffold:builder

	if err := mgr.AddHealthzCheck("healthz", healthz.Ping); err != nil {
		return err
	}

	if err := mgr.AddReadyzCheck("readyz", healthz.Ping); err != nil {
		return err
	}

	setupLog.Info("starting manager")

	if err := mgr.Start(ctrl.SetupSignalHandler()); err != nil {
		return err
	}

	return nil
}

func setupControllers(f *flags, mgr ctrl.Manager) error {
	err := controller.NewIngressController(mgr)
	if err != nil {
		return err
	}

	err = controller.NewMonitorController(mgr)
	if err != nil {
		return err
	}

	return nil
}
