# Extension Bootstrapping

An alternate approach to extension project creation.

## Concept

Templates are divided into two groups `templates/project` and `templates/shared`.
The former outlines the exact structure of an extension.
Every file in a project is run through the template engine and has access to a configuration object as well as all templates in `templates/shared`.

The most sophisticated example is the generation of the `extension.config.yml`.
The project template looks as follows:

```gotmpl
{{merge "shared/extension.config.yml" "shared/admin_extension/extension.config.yml"}}
```

Each referenced file is run through the template engine to substitute any potential placeholders.
The resulting fragments are then unmarshalled and deep merged.

## Experiment

To create the example project, simply run

```bash
go run .
```

It will create a `tmp/` directory that has the exact same structure as `templates/project/`.
