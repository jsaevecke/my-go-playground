package config

type HttpServer struct {
	Port string `env:"HTTP_SERVER_PORT" envDefault:"8080"`
}
