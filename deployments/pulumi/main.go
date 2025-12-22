package main

import (
	"fmt"

	// For the Migrator Job
	// For PodSpec, Containers, etc.

	// For ObjectMeta
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

func main() {
	pulumi.Run(func(ctx *pulumi.Context) error {
		k8sProvider, err := setupClusterProvider(ctx)
		if err != nil {
			return fmt.Errorf("error setting up cluster provider: %w", err)
		}

		namespaces, err := createNamespaces(ctx, k8sProvider)
		if err != nil {
			return fmt.Errorf("error on creating namespaces")
		}

		images, err := buildInternalImages(ctx)
		if err != nil {
			return fmt.Errorf("error on building images: %w", err)
		}

		dbDetails, err := deployDatabase(ctx, k8sProvider, namespaces)
		if err != nil {
			return fmt.Errorf("error on deploying database: %w", err)
		}

		if err := deployApplications(ctx, k8sProvider, images, namespaces, dbDetails); err != nil {
			return fmt.Errorf("error on deploying application: %w", err)
		}
		return nil
	})
}
