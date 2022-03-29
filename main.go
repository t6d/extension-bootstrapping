package main

import (
	"embed"
	"io/fs"
	"os"
	"strings"
)

//go:embed templates/*
var templates embed.FS

func main() {
	source_directory := "templates/project"
	target_directory := "tmp"

	fs.WalkDir(templates, source_directory, func(source string, d fs.DirEntry, err error) error {
		target := strings.Replace(source, source_directory, target_directory, 1)

		if d.IsDir() {
			if err := os.Mkdir(target, 0755); !os.IsExist(err) {
				panic(err)
			}
		} else {
			data, _ := templates.ReadFile(source)
			file, _ := os.Create(target)
			file.Write(data)
			file.Close()
		}
		return nil
	})
}
