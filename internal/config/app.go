package config

type App struct {
	AppEnvironment string `env:"APP_ENVIRONMENT" envDefault:"dev"`
	AppName        string `env:"APP_NAME" envDefault:"app"`
	AppLogLevel    string `env:"APP_LOG_LEVEL" envDefault:"info"`
}
