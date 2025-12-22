package config

type HttpServer struct {
	Port int `env:"HTTP_SERVER_PORT" envDefault:"8080"`
}
