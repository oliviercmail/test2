package {{ .package }}

{{ template "gocode/header-gentext.tpl" }}

{{ if .imports }}
import (
{{- range .imports }}
    {{ . }}
{{- end }}
)
{{- end }}

type (
	{{ range .groups }}
    {{ .struct }} struct {
      {{- range .options }}
        {{ .expIdent }} {{ .type }} `env:"{{ .env }}"`
      {{- end }}
    }
	{{ end }}
)

{{ range .groups }}
	// {{ .func }} initializes and returns a {{ .struct }} with default values
	//
	// This function is auto-generated
	func {{ .func }}() (o *{{ .struct }}) {
			o = &{{ .struct }}{
				{{- range  .options }}
					{{- if or .default }}
						{{ .expIdent }}: {{ .default }},
					{{- end }}
				{{- end }}
			}

			fill(o)

			// Function that allows access to custom logic inside the parent function.
			// The custom logic in the other file should be like:
			// func (o *{{ .struct }}) Defaults() {...}
			func(o interface{}) {
				if def, ok := o.(interface{ Defaults() }); ok {
					def.Defaults()
				}
			}(o)

			return
	}
{{ end }}
