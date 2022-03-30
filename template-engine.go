package main

import (
	"bytes"
	"io/fs"
	"strings"
	"text/template"

	"gopkg.in/yaml.v3"
)

func createTemplateEngine(shared fs.FS, config Config) *template.Template {
	t := &template.Template{}
	t.Funcs(helpers(t, config, shared))

	fs.WalkDir(shared, ".", func(path string, d fs.DirEntry, err error) error {
		if !d.IsDir() {
			data, _ := fs.ReadFile(shared, path)
			t.New("shared/" + path).Parse(string(data))
		}

		return nil
	})

	return t
}

func helpers(t *template.Template, config Config, shared fs.FS) template.FuncMap {
	return template.FuncMap{
		"merge": func(paths ...string) string {
			fragments := make([]fragment, 0, len(paths))

			for _, path := range paths {
				rawData, _ := fs.ReadFile(shared, strings.TrimPrefix(path, "shared/"))

				buffer := bytes.Buffer{}
				template.Must(t.New("").Parse(string(rawData))).Execute(&buffer, config)

				fragment := make(fragment)
				yaml.Unmarshal(buffer.Bytes(), fragment)
				fragments = append(fragments, fragment)
			}

			result := fragments[0]
			for _, fragment := range fragments[1:] {
				result = mergeFragments(result, fragment)
			}
			serializedResult, _ := yaml.Marshal(result)
			return string(serializedResult)
		},
	}
}

func mergeFragments(a, b fragment) fragment {
	out := make(fragment, len(a))
	for k, v := range a {
		out[k] = v
	}
	for k, v := range b {
		if v, ok := v.(fragment); ok {
			if bv, ok := out[k]; ok {
				if bv, ok := bv.(fragment); ok {
					out[k] = mergeFragments(bv, v)
					continue
				}
			}
		}
		out[k] = v
	}
	return out
}

type fragment = map[interface{}]interface{}
