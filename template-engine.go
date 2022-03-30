package main

import (
	"io/fs"
	"text/template"
)

func createTemplateEngine(shared fs.FS) *template.Template {
	t := template.Template{}
	t.Funcs(helpers)

	fs.WalkDir(shared, ".", func(path string, d fs.DirEntry, err error) error {
		if !d.IsDir() {
			data, _ := fs.ReadFile(shared, path)
			t.New("shared/" + path).Parse(string(data))
		}

		return nil
	})

	return &t
}

var helpers template.FuncMap = template.FuncMap{}
