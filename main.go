package main

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
)

func main() {
	if err := run(); err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "%s\n", err)
		os.Exit(1)
	}
}

func run() error {
	// set defaults
	// if kubeconfig absolute path not specified, read from KUBECONFIG, or load from home directory (see https://github.com/kubernetes/client-go/blob/master/examples/out-of-cluster-client-configuration/)
	var kubeconfig string
	if kubeconfig = os.Getenv("KUBECONFIG"); kubeconfig == "" {
		if home := homeDir(); home != "" {
			kubeconfig = filepath.Join(home, ".kube", "config")
		}
	}
	var kubeclient kubernetes.Interface
	var apiGroup string

	var cmd = &cobra.Command{
		Use:          "kapi-groups",
		Short:        "List all resources using a specified APIgroup",
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			var err error
			kubeclient, err = initializeKubeClient(kubeconfig)
			return err
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			return findResources(kubeclient, apiGroup)
		},
	}

	cmd.Flags().StringVar(&kubeconfig, "kubeconfig", kubeconfig, "(optional) Absolute path to the kubeconfig file")
	cmd.PersistentFlags().StringVar(&apiGroup, "api-group", apiGroup, "Name of the API Group from which you want to find resources")

	return cmd.Execute()
}

func initializeKubeClient(kubeconfig string) (*kubernetes.Clientset, error) {
	config, err := clientcmd.BuildConfigFromFlags("", kubeconfig)
	if err != nil {
		return nil, errors.Wrap(err, "error building kubeconfig")
	}
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		return nil, errors.Wrap(err, "error creating kube client")
	}
	return clientset, nil
}

func findResources(kubeClient kubernetes.Interface, apiGroup string) error {
	return nil
}

func homeDir() string {
	if h := os.Getenv("HOME"); h != "" {
		return h
	}
	return os.Getenv("USERPROFILE")
}
