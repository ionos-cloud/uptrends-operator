package main

import (
	"context"

	"github.com/ionos-cloud/uptrends-operator/signals"

	"github.com/ionos-cloud/uptrends-operator/monitor"
	"github.com/spf13/cobra"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/klog/v2"
	//+kubebuilder:scaffold:imports
)

type flags struct {
	KubeConfig string
	MasterURL  string
}

var f = &flags{}

var rootCmd = &cobra.Command{
	Use: "controller",
	RunE: func(cmd *cobra.Command, args []string) error {
		return run(cmd.Context())
	},
}

func init() {
	rootCmd.Flags().StringVar(&f.KubeConfig, "kubeconfig", f.KubeConfig, "Path to a kubeconfig. Only required if out-of-cluster.")
	rootCmd.Flags().StringVar(&f.MasterURL, "master", f.MasterURL, "The address of the Kubernetes API server. Overrides any value in kubeconfig. Only required if out-of-cluster.")
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		klog.Error(err, "unable to run controller")
	}
}

func run(ctx context.Context) error {
	klog.InitFlags(nil)

	stop := signals.SetupSignalHandler()

	cfg, err := clientcmd.BuildConfigFromFlags(f.MasterURL, f.KubeConfig)
	if err != nil {
		return err
	}

	clientset, err := kubernetes.NewForConfig(cfg)
	if err != nil {
		return err
	}

	informer := monitor.NewInformer(clientset, stop)
	controller := monitor.NewController(clientset, informer)

	err = controller.Run(stop)
	if err != nil {
		return err
	}

	return nil
}
