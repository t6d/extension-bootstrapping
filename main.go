package main

import (
	"embed"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
)

//go:embed templates/*
var templates embed.FS

func main() {
	config := Config{"test-project"}

	shared, err := fs.Sub(templates, "templates/shared")
	if err != nil {
		panic(err)
	}

	project, err := fs.Sub(templates, "templates/project")
	if err != nil {
		panic(err)
	}

	t := createTemplateEngine(shared)

	fs.WalkDir(project, ".", func(source string, d fs.DirEntry, err error) error {
		target := filepath.Join("tmp", source)
		fmt.Printf("Creating %s\n", target)

		if d.IsDir() {
			if err := os.Mkdir(target, 0755); err != nil && !os.IsExist(err) {
				panic(err)
			}
		} else {
			data, err := fs.ReadFile(project, source)
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

type Config struct {
	Name string
}
