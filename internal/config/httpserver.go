package config

type HttpServer struct {
	HttpServerPort int `env:"HTTP_SERVER_PORT" envDefault:"8080"`
}
