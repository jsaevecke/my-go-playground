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
