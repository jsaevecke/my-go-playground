package main

import (
	"fmt"

	"github.com/pulumi/pulumi-kubernetes/sdk/v4/go/kubernetes"
	"github.com/pulumi/pulumi-kubernetes/sdk/v4/go/kubernetes/apiextensions"
	"github.com/pulumi/pulumi-kubernetes/sdk/v4/go/kubernetes/helm/v3"
	metav1 "github.com/pulumi/pulumi-kubernetes/sdk/v4/go/kubernetes/meta/v1"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
	"github.com/pulumiverse/pulumi-time/sdk/go/time"
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

	wait, err := time.NewSleep(ctx, "wait-30-seconds", &time.SleepArgs{
		CreateDuration: pulumi.String("30s"),
	}, pulumi.Provider(p), pulumi.DependsOn([]pulumi.Resource{cnpgOperator}))
	if err != nil {
		return nil, fmt.Errorf("wait: %w", err)
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
	}, pulumi.Provider(p), pulumi.DependsOn([]pulumi.Resource{namespaces.dev, cnpgOperator, wait}))
	if err != nil {
		return nil, fmt.Errorf("failed to create cnpg cluster: %w", err)
	}

	// 2. Define your PG Cluster (simplified example)
	// In a real scenario, you'd use the CustomResource API
	return &DBDetails{
		Endpoint: pulumi.String("postgres-cluster-rw.default.svc.cluster.local"),
	}, nil
}
