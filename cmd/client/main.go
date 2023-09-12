package main

import (
	"log"

	"github.com/Nexadis/Storage/internal/client"
	"github.com/Nexadis/Storage/internal/client/config"
)

func main() {
	conf := config.New()
	c := client.New(conf)
	log.Fatal(c.Run())
}
