package config

import "flag"

type Config struct {
	Addr     string
	FileSave string
}

func New() *Config {
	c := &Config{}
	flag.StringVar(&c.Addr, "a", ":8080", "Address for HTTP server")
	flag.StringVar(&c.FileSave, "s", "transactions.log", "Filename for transactions file")
	flag.Parse()
	return c
}
