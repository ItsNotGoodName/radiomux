package config

import (
	"flag"
	"fmt"
	"net/url"

	"github.com/ItsNotGoodName/radiomux/internal/core"
	"github.com/caarlos0/env/v9"
	"github.com/rs/zerolog/log"
)

type Config struct {
	File       string   `env:"RADIOMUX_FILE"`
	HTTPHost   string   `env:"RADIOMUX_HTTP_HOST"`
	HTTPPort   int      `env:"RADIOMUX_HTTP_PORT"`
	HTTPURL    *url.URL `env:"-"`
	HTTPURLRaw string   `env:"RADIOMUX_HTTP_URL"`
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
	flags.StringVar(&c.HTTPURLRaw, "http-url", "", "HTTP public URL (e.g. http://127.0.0.1:8080).")
}

func (c *Config) Parse() error {
	err := env.Parse(c)
	if err != nil {
		return err
	}

	if c.HTTPURLRaw == "" {
		ips, err := core.PossiblePublicIPs()
		if err != nil {
			log.Err(err).Caller().Msg("Failed to list public ips")
		} else if len(ips) > 0 {
			ip := ips[0]
			c.HTTPURL, err = url.Parse(fmt.Sprintf("http://%s:%d", ip, c.HTTPPort))
			if err != nil {
				return err
			}
		}
	} else {
		c.HTTPURL, err = url.Parse(c.HTTPURLRaw)
		if err != nil {
			return err
		}
	}

	return nil
}
