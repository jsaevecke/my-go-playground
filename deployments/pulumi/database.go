package main

import (
	"fmt"

	"github.com/pulumi/pulumi-kubernetes/sdk/v4/go/kubernetes"
	"github.com/pulumi/pulumi-kubernetes/sdk/v4/go/kubernetes/apiextensions"
	appsv1 "github.com/pulumi/pulumi-kubernetes/sdk/v4/go/kubernetes/apps/v1" // For the Migrator Job

	// For PodSpec, Containers, etc.
	"github.com/pulumi/pulumi-kubernetes/sdk/v4/go/kubernetes/helm/v3"
	metav1 "github.com/pulumi/pulumi-kubernetes/sdk/v4/go/kubernetes/meta/v1" // For ObjectMeta
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

type DBDetails struct {
	Endpoint pulumi.StringInput
}

func deployDatabase(ctx *pulumi.Context, p *kubernetes.Provider, namespaces *Namespaces) (*DBDetails, error) {
	// setup cnpg operator in 'cnpg-system' namespace - required to spin up a cnpg cluster
	cnpgOperator, err := helm.NewChart(ctx, "cnpg-operator", helm.ChartArgs{
		Chart:     pulumi.String("cloudnative-pg"),
		Version:   pulumi.String("0.21.0"),
		Namespace: pulumi.String("cnpg-system"),
		FetchArgs: &helm.FetchArgs{
			Repo: pulumi.String("https://cloudnative-pg.github.io/charts"),
		},
	}, pulumi.Provider(p), pulumi.DependsOn([]pulumi.Resource{namespaces.cnpg}))
	if err != nil {
		return nil, fmt.Errorf("create cnpg-operator: %w", err)
	}

	operatorDeployment, err := appsv1.GetDeployment(ctx, "cnpg-operator-deployment",
		pulumi.ID("cnpg-system/cnpg-operator-cloudnative-pg"),
		nil,
		pulumi.Provider(p),
		pulumi.DependsOn([]pulumi.Resource{cnpgOperator}),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to get cnpg operator: %w", err)
	}

	_, err = apiextensions.NewCustomResource(ctx, "postgres-ha-cluster", &apiextensions.CustomResourceArgs{
		ApiVersion: pulumi.String("postgresql.cnpg.io/v1"),
		Kind:       pulumi.String("Cluster"),
		Metadata: &metav1.ObjectMetaArgs{
			Name:      pulumi.String("my-app-db"),
			Namespace: namespaces.dev.Metadata.Name(),
		},
		OtherFields: kubernetes.UntypedArgs{
			"spec": map[string]any{
				"instances": 2, // 1 Primary, 1 Replicas
				"storage": map[string]any{
					"size": "1Gi",
				},
				"bootstrap": map[string]any{
					"initdb": map[string]any{
						"database": "app_db",
						"owner":    "app_user",
					},
				},
			},
		},
	}, pulumi.Provider(p), pulumi.DependsOn([]pulumi.Resource{operatorDeployment, namespaces.dev}))
	if err != nil {
		return nil, fmt.Errorf("failed to create cnpg cluster: %w", err)
	}

	// 2. Define your PG Cluster (simplified example)
	// In a real scenario, you'd use the CustomResource API
	return &DBDetails{
		Endpoint: pulumi.String("postgres-cluster-rw.default.svc.cluster.local"),
	}, nil
}
