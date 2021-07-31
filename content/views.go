// +build !embed

package content

import (
	"html/template"
)

func embedStaticViews() {}

func parseTemplate(t *template.Template, lf layoutFile) *template.Template {
	var err error
	t, err = t.ParseFiles(lf.Files...)
	if err != nil {
		return nil
	}

	return t
}
