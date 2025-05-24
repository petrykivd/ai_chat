package config

type ServerConfig struct {
	Port  int  `envconfig:"SERVER_PORT" default:"8080"`
	Debug bool `envconfig:"SERVER_DEBUG" default:"false"`
}
