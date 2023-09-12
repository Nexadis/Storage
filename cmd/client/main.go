package main

import (
	"flag"
	"log"

	"github.com/Nexadis/Storage/internal/client"
	"github.com/Nexadis/Storage/internal/client/config"
)

func main() {
	conf := config.New()
	c := client.New(conf)
	err := c.Open()
	if err != nil {
		log.Fatal(err)
	}
	log.Fatal(c.DoCmd(flag.Args()))
}
