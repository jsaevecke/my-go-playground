package main

import (
	"github.com/pulumi/pulumi-kubernetes/sdk/v4/go/kubernetes"
	appsv1 "github.com/pulumi/pulumi-kubernetes/sdk/v4/go/kubernetes/apps/v1"
	batchv1 "github.com/pulumi/pulumi-kubernetes/sdk/v4/go/kubernetes/batch/v1" // For the Migrator Job
	corev1 "github.com/pulumi/pulumi-kubernetes/sdk/v4/go/kubernetes/core/v1"   // For PodSpec, Containers, etc.
	metav1 "github.com/pulumi/pulumi-kubernetes/sdk/v4/go/kubernetes/meta/v1"   // For ObjectMeta
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

func deployApplications(
	ctx *pulumi.Context,
	p *kubernetes.Provider,
	images *AppImages,
	namespaces *Namespaces,
	db *DBDetails,
) error {
	_, err := batchv1.NewJob(ctx, "migrator-job", &batchv1.JobArgs{
		Metadata: &metav1.ObjectMetaArgs{
			Name:      pulumi.String("db-migrator"),
			Namespace: namespaces.dev.Metadata.Name(),
		},
		Spec: &batchv1.JobSpecArgs{
			Template: &corev1.PodTemplateSpecArgs{
				Spec: &corev1.PodSpecArgs{
					Containers: corev1.ContainerArray{
						&corev1.ContainerArgs{
							Name:  pulumi.String("migrator"),
							Image: images.Migrator.Ref, // Use the image we just built
							Env: corev1.EnvVarArray{
								&corev1.EnvVarArgs{
									Name:  pulumi.String("DATABASE_PRIMARY_HOST"),
									Value: pulumi.String("postgres://user:pass@postgres-service:5432/db"),
								},
							},
						},
					},
					RestartPolicy: pulumi.String("OnFailure"),
				},
			},
		},
	})
	if err != nil {
		return err
	}

	appLabels := pulumi.StringMap{"app": pulumi.String("api")}
	_, err = appsv1.NewDeployment(ctx, "api-deploy", &appsv1.DeploymentArgs{
		Spec: &appsv1.DeploymentSpecArgs{
			Selector: &metav1.LabelSelectorArgs{MatchLabels: appLabels},
			Template: &corev1.PodTemplateSpecArgs{
				Metadata: &metav1.ObjectMetaArgs{Labels: appLabels, Namespace: namespaces.dev.Metadata.Name()},
				Spec: &corev1.PodSpecArgs{
					Containers: corev1.ContainerArray{
						&corev1.ContainerArgs{
							Name:  pulumi.String("api"),
							Image: images.API.Ref, // Use the image we just built
							Ports: corev1.ContainerPortArray{
								&corev1.ContainerPortArgs{ContainerPort: pulumi.Int(8080)},
							},
						},
					},
				},
			},
		},
	})
	if err != nil {
		return err
	}

	_, err = corev1.NewService(ctx, "api-service", &corev1.ServiceArgs{
		Metadata: &metav1.ObjectMetaArgs{Labels: appLabels, Namespace: namespaces.dev.Metadata.Name()},
		Spec: &corev1.ServiceSpecArgs{
			Ports: corev1.ServicePortArray{
				&corev1.ServicePortArgs{Port: pulumi.Int(80), TargetPort: pulumi.Int(8080)},
			},
			Selector: appLabels,
		},
	})
	return nil
}
