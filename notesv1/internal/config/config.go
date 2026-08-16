package config

type HTTPConfig struct {
	Port int `envconfig:"PORT" default:"8080"`
}

type DBConfig struct {
	DSN string `envconfig:"DSN" required:"true"`
}

type Config struct {
	HTTP HTTPConfig `envconfig:"HTTP"`
	DB   DBConfig   `envconfig:"DB"`
}
