package main

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/pulumi/pulumi-command/sdk/go/command/local"
	"github.com/pulumi/pulumi-kubernetes/sdk/v4/go/kubernetes"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

func setupClusterProvider(ctx *pulumi.Context) (*kubernetes.Provider, error) {
	cwd, _ := os.Getwd()
	kubeconfigPath := filepath.Join(cwd, "kind-kubeconfig")
	kindConfigPath := filepath.Join(cwd, "..", "kind-config.yaml")

	// create klind cluster with using the kind config and wait until it is setup
	cluster, err := local.NewCommand(ctx, "kind-cluster", &local.CommandArgs{
		// Wichtig: --kubeconfig ./kind-kubeconfig erstellt die Datei lokal im deployments/pulumi Ordner
		Create: pulumi.String(fmt.Sprintf("kind create cluster --name dev-cluster --config %s --kubeconfig %s --wait 5m", kindConfigPath, kubeconfigPath)),
		Delete: pulumi.String(fmt.Sprintf("kind delete cluster --name dev-cluster && rm -f %s", kubeconfigPath)),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create cluster: %w", err)
	}

	// try to get the kube config (the config might not exist on the first run or preview)
	kubeconfig := cluster.ID().ApplyT(func(_ pulumi.ID) (string, error) {
		if _, err := os.Stat(kubeconfigPath); os.IsNotExist(err) {
			// if it does not exist, we are going to ignore it for now (preview)
			return "", nil
		}

		time.Sleep(2 * time.Second)
		data, err := os.ReadFile(kubeconfigPath)
		if err != nil {
			return "", fmt.Errorf("failed to read kubeconfig at %s: %w", kubeconfigPath, err)
		}
		return string(data), nil
	}).(pulumi.StringOutput)

	k8sProvider, err := kubernetes.NewProvider(ctx, "k8s-provider", &kubernetes.ProviderArgs{
		Kubeconfig: kubeconfig,
	}, pulumi.DependsOn([]pulumi.Resource{cluster}))
	if err != nil {
		return nil, fmt.Errorf("failed to create new k8s provider: %w", err)
	}
	return k8sProvider, nil
}
