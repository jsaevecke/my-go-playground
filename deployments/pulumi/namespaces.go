package main

import (
	"fmt"

	"github.com/pulumi/pulumi-kubernetes/sdk/v4/go/kubernetes"
	corev1 "github.com/pulumi/pulumi-kubernetes/sdk/v4/go/kubernetes/core/v1"
	metav1 "github.com/pulumi/pulumi-kubernetes/sdk/v4/go/kubernetes/meta/v1"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

type Namespaces struct {
	cnpg *corev1.Namespace
	dev  *corev1.Namespace
}

func createNamespaces(ctx *pulumi.Context, p *kubernetes.Provider) (*Namespaces, error) {
	namespaceDev, err := corev1.NewNamespace(ctx, "dev-ns", &corev1.NamespaceArgs{
		Metadata: &metav1.ObjectMetaArgs{
			Name: pulumi.String("dev"),
		},
	}, pulumi.Provider(p))
	if err != nil {
		return nil, fmt.Errorf("failed to create 'dev' namespace: %w", err)
	}

	namespaceCnpg, err := corev1.NewNamespace(ctx, "cnpg-system-ns", &corev1.NamespaceArgs{
		Metadata: &metav1.ObjectMetaArgs{
			Name: pulumi.String("cnpg-system"),
		},
	}, pulumi.Provider(p))
	if err != nil {
		return nil, fmt.Errorf("create 'cnpg-system' namespace: %w", err)
	}

	return &Namespaces{
		dev:  namespaceDev,
		cnpg: namespaceCnpg,
	}, nil
}
