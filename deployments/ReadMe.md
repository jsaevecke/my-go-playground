```
.
├── cmd/
│   └── migrator/             # Application entry point
|       └── migration/
|           └── 00001_user_model.sql
|       └── .env
|       └── configuration.go
│       └── main.go           # Wire up everything (Dependency Injection)
├── deployments/
│   ├── kind-config.yaml       # Die Cluster-Konfiguration für Kind
│   └── pulumi/
│       ├── main.go            # Setups local kind cluster with cloudnative postgres HA cluster
│       ├── Pulumi.yaml        # Projekt-Metadaten
│       ├── Pulumi.dev.yaml    # Stack-spezifische Konfiguration
│       └── go.mod             # Go-Abhängigkeiten (Pulumi SDKs)
├── internal/               # Private application code
├── pkg/                    # Shared public utilities (optional)
├── go.mod
```

`cd deployments/pulumi`
`pulumi login --local`
`pulumi stack init dev`

```
# 1. Build the API image
docker build --build-arg APP_NAME=api -t my-go-playground/api:latest .

# 2. Build the Migrator image
docker build --build-arg APP_NAME=migrator -t my-go-playground/migrator:latest .

# 3. Side-load images into Kind
kind load docker-image my-go-playground/api:latest
kind load docker-image my-go-playground/migrator:latest

# 4. Deploy with Pulumi
cd deployments/pulumi && pulumi up
```
