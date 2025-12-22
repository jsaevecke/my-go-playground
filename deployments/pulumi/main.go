package main

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/pulumi/pulumi-command/sdk/go/command/local"
	"github.com/pulumi/pulumi-docker-build/sdk/go/dockerbuild"
	"github.com/pulumi/pulumi-kubernetes/sdk/v4/go/kubernetes"
	"github.com/pulumi/pulumi-kubernetes/sdk/v4/go/kubernetes/apiextensions"
	appsv1 "github.com/pulumi/pulumi-kubernetes/sdk/v4/go/kubernetes/apps/v1"
	batchv1 "github.com/pulumi/pulumi-kubernetes/sdk/v4/go/kubernetes/batch/v1" // For the Migrator Job
	corev1 "github.com/pulumi/pulumi-kubernetes/sdk/v4/go/kubernetes/core/v1"   // For PodSpec, Containers, etc.
	"github.com/pulumi/pulumi-kubernetes/sdk/v4/go/kubernetes/helm/v3"
	metav1 "github.com/pulumi/pulumi-kubernetes/sdk/v4/go/kubernetes/meta/v1" // For ObjectMeta
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

func main() {
	pulumi.Run(func(ctx *pulumi.Context) error {
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
			return fmt.Errorf("failed to create cluster: %w", err)
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
			return fmt.Errorf("failed to create new k8s provider: %w", err)
		}

		// setup cnpg operator in 'cnpg-system' namespace - required to spin up a cnpg cluster
		namespaceCNPG, err := corev1.NewNamespace(ctx, "cnpg-system-ns", &corev1.NamespaceArgs{
			Metadata: &metav1.ObjectMetaArgs{
				Name: pulumi.String("cnpg-system"),
			},
		}, pulumi.Provider(k8sProvider))
		if err != nil {
			return fmt.Errorf("failed to create 'cnpg-system' namespace: %w", err)
		}

		cnpgOperator, err := helm.NewChart(ctx, "cnpg-operator", helm.ChartArgs{
			Chart:     pulumi.String("cloudnative-pg"),
			Version:   pulumi.String("0.21.0"),
			Namespace: pulumi.String("cnpg-system"),
			FetchArgs: &helm.FetchArgs{
				Repo: pulumi.String("https://cloudnative-pg.github.io/charts"),
			},
		}, pulumi.Provider(k8sProvider), pulumi.DependsOn([]pulumi.Resource{namespaceCNPG}))
		if err != nil {
			return err
		}

		operatorDeployment, err := appsv1.GetDeployment(ctx, "cnpg-operator-deployment",
			pulumi.ID("cnpg-system/cnpg-operator-cloudnative-pg"),
			nil,
			pulumi.Provider(k8sProvider),
			pulumi.DependsOn([]pulumi.Resource{cnpgOperator}),
		)
		if err != nil {
			return fmt.Errorf("failed to get cnpg operator: %w", err)
		}

		// setup the cnpg cluster
		namespaceDev, err := corev1.NewNamespace(ctx, "dev-ns", &corev1.NamespaceArgs{
			Metadata: &metav1.ObjectMetaArgs{
				Name: pulumi.String("dev"),
			},
		}, pulumi.Provider(k8sProvider))
		if err != nil {
			return fmt.Errorf("failed to create 'dev' namespace: %w", err)
		}

		_, err = apiextensions.NewCustomResource(ctx, "postgres-ha-cluster", &apiextensions.CustomResourceArgs{
			ApiVersion: pulumi.String("postgresql.cnpg.io/v1"),
			Kind:       pulumi.String("Cluster"),
			Metadata: &metav1.ObjectMetaArgs{
				Name:      pulumi.String("my-app-db"),
				Namespace: namespaceDev.Metadata.Name(),
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
		}, pulumi.Provider(k8sProvider), pulumi.DependsOn([]pulumi.Resource{operatorDeployment, namespaceDev}))
		if err != nil {
			return fmt.Errorf("failed to create cnpg cluster: %w", err)
		}

		apiImage, err := dockerbuild.NewImage(ctx, "api-image", &dockerbuild.ImageArgs{
			Context: &dockerbuild.BuildContextArgs{
				Location: pulumi.String("../../"),
			},
			Dockerfile: &dockerbuild.DockerfileArgs{
				Location: pulumi.String("../../Dockerfile"),
			},
			BuildArgs: pulumi.StringMap{
				"APP_NAME": pulumi.String("api"),
			},
			Tags: pulumi.StringArray{pulumi.String("my-go-playground/api:latest")},
			Push: pulumi.Bool(false), // Don't push to registry, we use local kind
		})
		if err != nil {
			return err
		}

		migratorImage, err := dockerbuild.NewImage(ctx, "migrator-image", &dockerbuild.ImageArgs{
			Context: &dockerbuild.BuildContextArgs{
				Location: pulumi.String("../../"),
			},
			Dockerfile: &dockerbuild.DockerfileArgs{
				Location: pulumi.String("../../Dockerfile"),
			},
			BuildArgs: pulumi.StringMap{
				"APP_NAME": pulumi.String("migrator"),
			},
			Tags: pulumi.StringArray{pulumi.String("my-go-playground/migrator:latest")},
			Push: pulumi.Bool(false),
		})
		if err != nil {
			return err
		}

		_, err = batchv1.NewJob(ctx, "migrator-job", &batchv1.JobArgs{
			Metadata: &metav1.ObjectMetaArgs{
				Name: pulumi.String("db-migrator"),
			},
			Spec: &batchv1.JobSpecArgs{
				Template: &corev1.PodTemplateSpecArgs{
					Spec: &corev1.PodSpecArgs{
						Containers: corev1.ContainerArray{
							&corev1.ContainerArgs{
								Name:  pulumi.String("migrator"),
								Image: migratorImage.Ref, // Use the image we just built
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
					Metadata: &metav1.ObjectMetaArgs{Labels: appLabels},
					Spec: &corev1.PodSpecArgs{
						Containers: corev1.ContainerArray{
							&corev1.ContainerArgs{
								Name:  pulumi.String("api"),
								Image: apiImage.Ref, // Use the image we just built
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
			Metadata: &metav1.ObjectMetaArgs{Labels: appLabels},
			Spec: &corev1.ServiceSpecArgs{
				Ports: corev1.ServicePortArray{
					&corev1.ServicePortArgs{Port: pulumi.Int(80), TargetPort: pulumi.Int(8080)},
				},
				Selector: appLabels,
			},
		})
		return nil
	})
}
