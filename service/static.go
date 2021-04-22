// +build !embed

package service

import (
	"log"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/favicon"
)

func embedStaticResource(app *fiber.App) {
	app.Static("/", "static")
	log.Println("set static dir:", "static")
}

func setFavicon(app *fiber.App) func(*fiber.Ctx) error {
	const path = "static/image/favicon.ico"
	if _, err := os.Stat(path); err != nil {
		return func(*fiber.Ctx) error {
			log.Println("ignore favicon.ico")
			return nil
		}
	}

	return favicon.New(favicon.Config{
		File: path,
	})
}
