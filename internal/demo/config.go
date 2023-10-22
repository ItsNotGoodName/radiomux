package demo

import (
	"flag"
	"fmt"

	"github.com/caarlos0/env/v9"
)

type Config struct {
	Port   int    `env:"PORT"`
	Domain string `env:"DOMAIN"`
	Host   string `env:"-"`
}

func (c *Config) HTTPURL() string {
	return fmt.Sprintf("https://%s", c.Domain)
}

var Default = struct {
	Port   int
	Domain string
}{
	Port: 8080,
}

func New() *Config {
	return &Config{
		Port: Default.Port,
	}
}

func (c *Config) WithFlag(flags *flag.FlagSet) {
	flags.IntVar(&c.Port, "port", Default.Port, "Port to listen on.")
	flags.StringVar(&c.Domain, "domain", Default.Domain, "Domain for application (e.g. 192.168.1.100, radiomux.example.com:8080).")
}

func (c *Config) Parse() error {
	return env.Parse(c)
}

