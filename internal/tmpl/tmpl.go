package tmpl

import (
	"html/template"
	"os"
)

func RenderTemplate(name, tmpl, filePath string, data any) error {
	t, err := template.New(name).Parse(tmpl)
	if err != nil {
		return err
	}
	f, err := os.Create(filePath)
	if err != nil {
		return err
	}
	defer f.Close()
	if err := t.Execute(f, data); err != nil {
		return err
	}
	return nil
}
