package config

import "flag"

type Config struct {
	Addr     string
	GRPC     string
	FileSave string
	DBURI    string
}

func New() *Config {
	c := &Config{}
	flag.StringVar(&c.Addr, "a", ":8080", "Address for HTTP server")
	flag.StringVar(&c.GRPC, "g", ":50051", "Address for gRPC server")
	flag.StringVar(&c.FileSave, "s", "transactions.log", "Filename for transactions file")
	flag.StringVar(&c.DBURI, "d", "", "URI for db connect")
	flag.Parse()
	return c
}
