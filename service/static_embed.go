// +build embed

package service

import (
	"strconv"

	"github.com/gobuffalo/packr/v2"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/filesystem"
)

var (
	staticbox = packr.New("WebStatic", "../static")
)

func embedStaticResource(app *fiber.App) {
	// 将静态资源打包，便于发布
	app.Use("/", filesystem.New(filesystem.Config{
		Root: staticbox,
	}))
}

func setFavicon(app *fiber.App) func(c *fiber.Ctx) error {
	const (
		hType         = "image/x-icon"
		hAllow        = "GET, HEAD, OPTIONS"
		hZero         = "0"
		hCacheControl = "public, max-age=31536000"
	)

	var (
		err     error
		icon    []byte
		iconLen string
	)

	const path = "image/favicon.ico"
	icon, err = staticbox.Find(path)
	if err != nil {
		panic(err)
	}

	iconLen = strconv.Itoa(len(icon))

	return func(c *fiber.Ctx) error {
		// Only respond to favicon requests
		if len(c.Path()) != 12 || c.Path() != "/favicon.ico" {
			return c.Next()
		}

		// Only allow GET, HEAD and OPTIONS requests
		if c.Method() != fiber.MethodGet && c.Method() != fiber.MethodHead {
			if c.Method() != fiber.MethodOptions {
				c.Status(fiber.StatusMethodNotAllowed)
			} else {
				c.Status(fiber.StatusOK)
			}
			c.Set(fiber.HeaderAllow, hAllow)
			c.Set(fiber.HeaderContentLength, hZero)
			return nil
		}

		// Serve cached favicon
		if len(icon) > 0 {
			c.Set(fiber.HeaderContentLength, iconLen)
			c.Set(fiber.HeaderContentType, hType)
			c.Set(fiber.HeaderCacheControl, hCacheControl)
			return c.Status(fiber.StatusOK).Send(icon)
		}

		return c.SendStatus(fiber.StatusNoContent)
	}
}
