package main

import (
	"log"

	"github.com/goSeeFuture/gblogconfigs"
	"github.com/goSeeFuture/gblogcontent"
	"github.com/goSeeFuture/gblogservice"

	"github.com/gofiber/fiber/v2"
)

func main() {
	configs.Load()
	content.Load()

	app := fiber.New()

	service.Mux(app)
	err := app.Listen(configs.Setting.Listen)
	if err != nil {
		log.Fatal(err)
	}
}
