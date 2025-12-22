package config

type App struct {
	Environment string `env:"APP_ENVIRONMENT" envDefault:"dev"`
	Name        string `env:"APP_NAME" envDefault:"app"`
	LogLevel    string `env:"APP_LOG_LEVEL" envDefault:"info"`
}
