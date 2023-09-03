package config

import "flag"

type Config struct {
	Addr string
}

func New() *Config {
	c := &Config{}
	flag.StringVar(&c.Addr, "a", ":8080", "Address for HTTP server")
	flag.Parse()
	return c
}
