package config

import (
	"flag"

	"github.com/caarlos0/env/v9"
)

type Config struct {
	File     string `env:"RADIOMUX_FILE"`
	HTTPHost string `env:"RADIOMUX_HTTP_HOST"`
	HTTPPort int    `env:"RADIOMUX_HTTP_PORT"`
	HTTPURL  string `env:"RADIOMUX_HTTP_URL"`
}

var Default = struct {
	File     string
	HTTPPort int
}{
	File:     "radiomux.json",
	HTTPPort: 8080,
}

func New() *Config {
	return &Config{
		File:     Default.File,
		HTTPPort: Default.HTTPPort,
	}
}

func (c *Config) WithFlag(flags *flag.FlagSet) {
	flags.StringVar(&c.File, "file", Default.File, "File path to JSON database.")

	flags.StringVar(&c.HTTPHost, "http-host", "", "HTTP host to listen on.")
	flags.IntVar(&c.HTTPPort, "http-port", Default.HTTPPort, "HTTP port to listen on.")
	flags.StringVar(&c.HTTPURL, "http-url", "", "HTTP public URL (e.g. http://127.0.0.1:8080).")
}

func (c *Config) Parse() error {
	return env.Parse(c)
}
