package main

import (
	"bytes"
	"io/fs"
	"strings"
	"text/template"
)

func createTemplateEngine(shared fs.FS) *template.Template {
	t := &template.Template{}
	t.Funcs(helpers(t, shared))

	fs.WalkDir(shared, ".", func(path string, d fs.DirEntry, err error) error {
		if !d.IsDir() {
			data, _ := fs.ReadFile(shared, path)
			t.New("shared/" + path).Parse(string(data))
		}

		return nil
	})

	return t
}

func helpers(t *template.Template, shared fs.FS) template.FuncMap {
	return template.FuncMap{
		"merge": func(paths ...string) string {
			buffer := bytes.Buffer{}

			for _, path := range paths {
				data, _ := fs.ReadFile(shared, strings.TrimPrefix(path, "shared/"))
				buffer.Write(data)
			}

			builder := strings.Builder{}
			template.Must(t.New("").Parse(string(buffer.Bytes()))).Execute(&builder, Config{"Foo"})
			return builder.String()
		},
	}
}
