package main

import (
	"fmt"

	"github.com/pulumi/pulumi-command/sdk/go/command/local"
	"github.com/pulumi/pulumi-docker-build/sdk/go/dockerbuild"
	"github.com/pulumi/pulumi-kubernetes/sdk/v4/go/kubernetes" // For the Migrator Job

	// For PodSpec, Containers, etc.
	// For ObjectMeta
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

func deployApplications(
	ctx *pulumi.Context,
	p *kubernetes.Provider,
	namespaces *Namespaces,
	db *DBDetails,
	clusterDetails *ClusterDetails,
) error {
	var (
		applications = []string{"api", "migrator"}
	)
	for _, app := range applications {
		img, err := dockerbuild.NewImage(ctx, app+"-image", &dockerbuild.ImageArgs{
			Context:    &dockerbuild.BuildContextArgs{Location: pulumi.String("../../")},
			Dockerfile: &dockerbuild.DockerfileArgs{Location: pulumi.String("../../Dockerfile")},
			BuildArgs:  pulumi.StringMap{"APP_NAME": pulumi.String(app)},
			Tags:       pulumi.StringArray{pulumi.String("my-go-playground/" + app)},
			Push:       pulumi.Bool(false),
			Exports: dockerbuild.ExportArray{&dockerbuild.ExportArgs{
				Docker: &dockerbuild.ExportDockerArgs{},
			}},
		})
		if err != nil {
			return fmt.Errorf("new %s image: %w", app, err)
		}

		// Load API Image
		_, err = local.NewCommand(ctx, fmt.Sprintf("load-%s-to-kind", app), &local.CommandArgs{
			Create: pulumi.String(fmt.Sprintf(
				"kind load docker-image --name %s my-go-playground/%s",
				clusterDetails.Name, app,
			)),
			// The Triggers field ensures this re-runs if the image changes.
			// We use the image's ID or Digest so Pulumi knows when it's "new".
			Triggers: pulumi.Array{img.Ref},
		}, pulumi.Provider(p), pulumi.DependsOn([]pulumi.Resource{img, clusterDetails.Ref}))
		if err != nil {
			return fmt.Errorf("load api image to kind: %w", err)
		}
	}

	/*
		_, err = batchv1.NewJob(ctx, "migrator-job", &batchv1.JobArgs{
			Metadata: &metav1.ObjectMetaArgs{
				Name:      pulumi.String("db-migrator"),
				Namespace: namespaces.dev.Metadata.Name(),
			},
			Spec: &batchv1.JobSpecArgs{
				Template: &corev1.PodTemplateSpecArgs{
					Spec: &corev1.PodSpecArgs{
						Containers: corev1.ContainerArray{&corev1.ContainerArgs{
							Name:            pulumi.String("migrator"),
							Image:           migratorImage.Ref, // Use the image we just built
							ImagePullPolicy: pulumi.String("IfNotPresent"),
							Env: corev1.EnvVarArray{&corev1.EnvVarArgs{
								Name:  pulumi.String("DATABASE_PRIMARY_HOST"),
								Value: pulumi.String("postgres://user:pass@postgres-service:5432/db"),
							}}},
						},
						RestartPolicy: pulumi.String("OnFailure"),
					},
				},
			},
		}, pulumi.Provider(p), pulumi.DependsOn([]pulumi.Resource{namespaces.dev, loadedMigrator}))
		if err != nil {
			return fmt.Errorf("create migrator batch job: %w", err)
		}

		appLabels := pulumi.StringMap{"app": pulumi.String("api")}
		apiDeployment, err := appsv1.NewDeployment(ctx, "api-deploy", &appsv1.DeploymentArgs{
			Spec: &appsv1.DeploymentSpecArgs{
				Selector: &metav1.LabelSelectorArgs{MatchLabels: appLabels},
				Template: &corev1.PodTemplateSpecArgs{
					Metadata: &metav1.ObjectMetaArgs{Labels: appLabels, Namespace: namespaces.dev.Metadata.Name()},
					Spec: &corev1.PodSpecArgs{
						Containers: corev1.ContainerArray{
							&corev1.ContainerArgs{
								Name:            pulumi.String("api"),
								Image:           apiImage.Ref, // Use the image we just built
								ImagePullPolicy: pulumi.String("IfNotPresent"),
								Ports: corev1.ContainerPortArray{
									&corev1.ContainerPortArgs{ContainerPort: pulumi.Int(8080)},
								},
							},
						},
					},
				},
			},
		}, pulumi.Provider(p), pulumi.DependsOn([]pulumi.Resource{namespaces.dev, loadedAPI}))
		if err != nil {
			return fmt.Errorf("create api deployment: %w", err)
		}

		_, err = corev1.NewService(ctx, "api-service", &corev1.ServiceArgs{
			Metadata: &metav1.ObjectMetaArgs{Labels: appLabels, Namespace: namespaces.dev.Metadata.Name()},
			Spec: &corev1.ServiceSpecArgs{
				Ports: corev1.ServicePortArray{
					&corev1.ServicePortArgs{Port: pulumi.Int(80), TargetPort: pulumi.Int(8080)},
				},
				Selector: appLabels,
			},
		}, pulumi.Provider(p), pulumi.DependsOn([]pulumi.Resource{namespaces.dev, loadedAPI, apiDeployment}))
		if err != nil {
			return fmt.Errorf("new api service: %w", err)
		}
	*/
	return nil
}
