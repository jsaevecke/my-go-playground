# Infrastructure as Code (Pulumi)

This directory contains the modular Pulumi program used to provision the Kubernetes environment, build container images, and manage the database lifecycle.

## 📂 Modular Structure

- `main.go`: Orchestrates the sequence of events.
- `cluster.go`: Configures the Kubernetes Provider (Kind).
- `images.go`: Handles Docker builds for `api` and `migrator`.
- `database.go`: Provisions the CloudNativePG (CNPG) operator and clusters.
- `apps.go`: Defines the K8s Deployments, Services, and Migration Jobs.

## ⚙️ Configuration

Configuration is managed via Pulumi Stacks (`Pulumi.<stack>.yaml`).

### Adding New Settings

To add a new environment variable for an application:

1. **Plain Text:** `pulumi config set my-key my-value`
2. **Secret:** `pulumi config set my-secret my-value --secret`

The value will be added to your current stack file and can be accessed in `main.go` using the `config` package.

### Common Stacks

- `dev`: Local development targeting the `kind-my-cluster` context.

## 🔄 Deployment Lifecycle

1. **Image Building**: Pulumi uses the `docker-build` provider to build binaries into images using the root `Dockerfile`.
2. **DB Operator**: CNPG is installed via Helm to provide High-Availability Postgres.
3. **Migration Job**: A K8s `Job` is triggered. The API deployment depends on this job completing successfully.
4. **App Deployment**: The Go API is deployed with automated readiness/liveness probes.

## 🛠 Adjusting Resources

To change CPU/Memory limits or replica counts, navigate to `apps.go` and modify the `ContainerArgs` within the Deployment resource.
