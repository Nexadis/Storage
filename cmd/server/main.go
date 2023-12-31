package main

import (
	"log"

	"github.com/Nexadis/Storage/internal/server"
	"github.com/Nexadis/Storage/internal/server/config"
)

func main() {
	c := config.New()
	s := server.New(c)
	log.Fatal(s.Run())
}
