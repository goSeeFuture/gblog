package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

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
	go func() {
		err := app.Listen(configs.Setting.Listen)
		if err != nil {
			log.Fatal(err)
		}
	}()

	service.PrintAuthorSecret()

	// 等待关服信号，如 Ctrl+C、kill -2、kill -3、kill -15
	ch := make(chan os.Signal, 1)
	signal.Notify(ch, syscall.SIGTERM, syscall.SIGQUIT, syscall.SIGINT)
	<-ch
}
