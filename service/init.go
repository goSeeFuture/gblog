package service

import (
	"errors"
	"gblog/configs"
	"log"

	"github.com/gofiber/fiber/v2"

	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/gofiber/fiber/v2/middleware/session"
)

var (
	store = session.New()
)

func Mux(app *fiber.App) {
	setupMiddleware(app)

	app.Get("/", homePage)
	app.Get("/list", listPage)

	app.Get("/list/:page", listPage)
	app.Get("/articles/*", articlePage)
	app.Get("/tag/:tag/:page?", tagPage)
	app.Get("/category/:categoryid/:page?", categoryPage)
}

func setupMiddleware(app *fiber.App) {
	embedStaticResource(app)
	setArticleReference(app)
	app.Use(recover.New())
	app.Use(setFavicon(app))
}

func pageNumber(c *fiber.Ctx, pageKey string) int {
	sess, err := store.Get(c)
	if err != nil {
		panic(err) // middleware catch panic
	}

	pn := sess.Get(pageKey)
	if pn == nil {
		pn = 1
		sess.Set(pageKey, 1)
		if err := sess.Save(); err != nil {
			panic(err)
		}
	}

	curPage, ok := pn.(int)
	if !ok {
		panic(errors.New("invalid page number"))
	}
	return curPage
}

func setArticleReference(app *fiber.App) {
	if configs.Setting.ArticleReferenceDir != "" {
		app.Static("/", configs.Setting.ArticleReferenceDir)
		log.Println("set static dir:", configs.Setting.ArticleReferenceDir)
	}

	for _, dir := range configs.Setting.ArticleReferenceDirs {
		app.Static("/", dir)
		log.Println("set static dir:", dir)
	}
}
