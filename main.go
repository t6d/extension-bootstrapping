package main

import (
	"embed"
	"io/fs"
	"os"
	"path"
	"strings"
	"text/template"
)

//go:embed templates/*
var templates embed.FS

func main() {
	config := Config{"test-project"}

	sourceDirectory := "templates/project"
	targetDirectory := "tmp"

	t := createTemplateEngine()

	fs.WalkDir(templates, sourceDirectory, func(source string, d fs.DirEntry, err error) error {
		target := strings.Replace(source, sourceDirectory, targetDirectory, 1)

		if d.IsDir() {
			if err := os.Mkdir(target, 0755); !os.IsExist(err) {
				panic(err)
			}
		} else {
			data, err := templates.ReadFile(source)
			if err != nil {
				panic(err)
			}

			t.New(source)
			if _, err := t.Parse(string(data)); err != nil {
				panic(err)
			}

			file, err := os.Create(target)
			if err != nil {
				panic(err)
			}
			defer file.Close()

			if err := t.Execute(file, config); err != nil {
				panic(err)
			}
		}
		return nil
	})
}

func createTemplateEngine() (t template.Template) {
	if _, err := t.ParseFS(templates, "templates/data/*"); err != nil {
		panic(err)
	}
	t.Funcs(helpers)

	return
}

type Config struct {
	Name string
}

var helpers template.FuncMap = template.FuncMap{
	"render": func(templatePath string) string {
		data, err := templates.ReadFile(path.Join("templates", templatePath))
		if err != nil {
			panic(err)
		}
		return string(data)
	},
}
