# My Go Playground

A production-grade implementation of **Hexagonal Architecture** (Ports & Adapters) in Go, featuring automated infrastructure via Pulumi and Kubernetes (Kind).

## 🏗 Architecture

The project follows the Hexagonal Architecture pattern to ensure business logic remains independent of infrastructure, frameworks, and tools.

- **Internal/Domain**: Pure business logic and Interface (Port) definitions. Zero external dependencies.
- **Internal/Adapter**: Implementations (Adapters) for SQL (Goose), HTTP, and other external services.
- **Internal/Infrastructure**: Shared utilities like logging, configuration, and database connection pooling.
- **Cmd/**: Application entry points (Bootstrap/Composition Root).

## 🚀 Getting Started

### Prerequisites

- [Go 1.23+](https://go.dev/)
- [Docker](https://www.docker.com/)
- [Kind](https://kind.sigs.k8s.io/)
- [Task](https://taskfile.dev/)
- [Pulumi CLI](https://www.pulumi.com/docs/install/)

### Local Development Loop

The project uses `Taskfile.yaml` to orchestrate all development steps.

1. **Initialize Cluster:**
   ```bash
   task init
   ```
2. **Build & Deploy:**

   ```bash
   task up
   ```

3. **View Logs:**
   ```bash
   task logs:api
   ```

## 🛠 Best Practices

Dependency Inversion: High-level domain logic does not import low-level adapters.
Configuration: Scoped configuration—apps only see the environment variables they require.
Security: Pre-commit hooks via Gitleaks prevent secrets from entering Git history.
Slim Images: Multi-stage Docker builds using distroless/static for minimal attack surface.
