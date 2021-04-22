package content

import (
	"errors"
	"html/template"
	"io/ioutil"
	"log"
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

func Layout() error {
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

func Footer() template.HTML {
	footer, err := ioutil.ReadFile(configs.Setting.WebsiteFooter)
	if err != nil {
		log.Println("load website footer failed:", err)
		return ""
	}

	data := markdown2HTML(footer)
	if data == nil {
		log.Println("parse website footer failed")
		return ""
	}

	return template.HTML(data)
}
