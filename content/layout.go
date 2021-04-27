package content

import (
	"errors"
	"html/template"
	"io/ioutil"
	"log"
	"path/filepath"
	"strconv"
	"sync"
	"time"

	configs "github.com/goSeeFuture/gblog/configs"

	"github.com/gofiber/fiber/v2"
	"github.com/valyala/bytebufferpool"
)

type layoutFile struct {
	Name     string
	Files    []string
	Template *template.Template
}

var (
	layout      sync.Map
	layoutfiles = []layoutFile{
		{Name: "404", Files: []string{"views/index.html", "views/header.html", "views/footer.html", "views/mathjax.html", "views/404.html"}},
		{Name: "list", Files: []string{"views/index.html", "views/header.html", "views/footer.html", "views/mathjax.html", "views/list.html"}},
		{Name: "article", Files: []string{"views/index.html", "views/header.html", "views/footer.html", "views/mathjax.html", "views/article.html"}},
	}
)

func initLayoutTemplate() error {
	// 仅在go build -tags=embed有效
	embedStaticViews()

	for _, e := range layoutfiles {
		t := template.New(e.Name).Funcs(template.FuncMap{
			"IsDigit":    isDigit,
			"FormatTime": fmtTime,
		})
		t = parseTemplate(t, e)
		if t == nil {
			return errors.New("parse view template failed")
		}

		e.Template = t
		layout.Store(e.Name, e)
	}

	return nil
}

func Render(c *fiber.Ctx, tpl string, data interface{}) error {
	v, exist := layout.Load(tpl)
	if !exist {
		return errors.New("no layout")
	}

	buf := bytebufferpool.Get()
	defer bytebufferpool.Put(buf)

	err := v.(layoutFile).Template.ExecuteTemplate(buf, "index.html", data)
	if err != nil {
		return err
	}

	c.Response().Header.SetContentType(fiber.MIMETextHTMLCharsetUTF8)
	c.Context().SetBody(buf.Bytes())
	return c.SendStatus(fiber.StatusOK)
}

func isDigit(s string) bool {
	_, err := strconv.Atoi(s)
	return err == nil
}

func fmtTime(t time.Time) string {
	if t.IsZero() {
		return "-"
	}
	return t.Format(configs.TimeLayout)
}

func Page404() template.HTML {
	if !configs.Setting.CustomWebsite404 {
		return template.HTML("")
	}

	p404, err := ioutil.ReadFile(filepath.Join(configs.Setting.AbsArticleDir, configs.CustomPage404))
	if err != nil {
		log.Println("load website 404 page failed:", err)
		return template.HTML("")
	}

	data := markdown2HTML(p404)
	if data == nil {
		log.Println("parse website 404 page failed")
		return template.HTML("")
	}

	return template.HTML(data)
}

func Footer() template.HTML {
	if !configs.Setting.CustomWebsiteFooter {
		return template.HTML("")
	}

	footer, err := ioutil.ReadFile(filepath.Join(configs.Setting.AbsArticleDir, configs.CustomPageFooter))
	if err != nil {
		log.Println("load website footer failed:", err)
		return template.HTML("")
	}

	data := markdown2HTML(footer)
	if data == nil {
		log.Println("parse website footer failed")
		return template.HTML("")
	}

	return template.HTML(data)
}
