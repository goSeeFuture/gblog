// +build !embed

package content

import (
	"html/template"
)

func embedStaticViews() {}

func parseTemplate(t *template.Template, lf layoutFile) (*template.Template, error) {
	var err error
	t, err = t.ParseFiles(lf.Files...)
	if err != nil {
		return nil, err
	}

	return t, nil
}
