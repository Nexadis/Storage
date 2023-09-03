package main

import (
	"log"

	"github.com/Nexadis/Storage/internal/config"
	"github.com/Nexadis/Storage/internal/server"
)

func main() {
	c := config.New()
	s := server.New(c)
	log.Fatal(s.Run())
}
