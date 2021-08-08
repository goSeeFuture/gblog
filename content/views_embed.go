// +build embed

package content

import (
	"html/template"
	"log"
	"path/filepath"
	"strings"

	"github.com/gobuffalo/packr/v2"
)

var (
	// 视图文件
	views *packr.Box
)

func embedStaticViews() {
	views = packr.New("ViewTemplates", "../views")
}

func parseTemplate(t *template.Template, lf layoutFile) (*template.Template, error) {
	var err error
	var tpl string
	for _, f := range lf.Files {
		f = strings.TrimPrefix(f, "views/")
		tpl, err = views.FindString(f)
		if err != nil {
			log.Println("not found static template:", f)
			return nil, err
		}

		var tmpl *template.Template
		name := filepath.Base(f)
		if t == nil {
			t = template.New(name)
		}
		if name == t.Name() {
			tmpl = t
		} else {
			tmpl = t.New(name)
		}

		_, err = tmpl.Parse(tpl)
		if err != nil {
			log.Println("parse template", f, " failed:", err)
			return nil, err
		}
	}

	return t, nil
}
