package main

import (
	"log"

	"github.com/goSeeFuture/gblog/configs"
	"github.com/goSeeFuture/gblog/content"
	"github.com/goSeeFuture/gblog/service"

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
