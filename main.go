package main

import (
	"embed"
	"io/fs"
	"os"
	"strings"
	"text/template"
)

//go:embed templates/*
var templates embed.FS

func main() {
	config := Config{"test-project"}

	source_directory := "templates/project"
	target_directory := "tmp"

	fs.WalkDir(templates, source_directory, func(source string, d fs.DirEntry, err error) error {
		target := strings.Replace(source, source_directory, target_directory, 1)

		if d.IsDir() {
			if err := os.Mkdir(target, 0755); !os.IsExist(err) {
				panic(err)
			}
		} else {
			data, err := templates.ReadFile(source)
			if err != nil {
				panic(err)
			}

			t := template.New(source)
			t.Funcs(helpers)

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

type Config struct {
	Name string
}

var helpers template.FuncMap = template.FuncMap{
	"render": func(path string) string {
		return path
	},
}
