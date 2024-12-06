package proto

import (
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
	"strings"
	"text/template"
)

var fileTemplate = `syntax = "{{ .Syntax }}";

package {{ replace .Package "-" "." }};

{{- range .Options }}
option {{ .Name }} = "{{ .Value }}";
{{- end }}

{{- range .Messages }}
{{range .Descs }}
// {{ . }}
{{- end }}
{{- $fieldCount := len .Fields }}
message {{ title .Name }} { {{- if eq $fieldCount 0 -}} }{{ else }}
{{- range $index,$value := .Fields }}
{{- range .Descs }}
    // {{ . }}
{{- end }}
    {{ if .Repeated }}repeated {{ end }}{{ .TypeName }} {{ .Name }} = {{ add $index 1 }};
{{- end }}
}
{{- end }}
{{- end }}
{{ range .Service.Descs }}
// {{ . }}
{{- end }}
service {{ title .Service.Name }} {
{{- range .Rpcs }}
{{- range .Descs }}
    // {{ . }}
{{- end }}
    rpc {{ title .Name }} ({{ if ne .Request nil }}{{ .Request.Name }}{{ else }}Empty{{ end }}) returns ({{ if ne .Response nil }}{{ .Response.Name }}{{ else }}Empty{{ end }});
{{- end }}
}
`

var funcMap = template.FuncMap{
	"add": func(x, y int) int {
		return x + y
	},
	"title": func() func(string) string {
		caser := cases.Title(language.English, cases.NoLower)
		return func(s string) string {
			ss := strings.Split(s, "-")
			for i := range ss {
				ss[i] = caser.String(ss[i])
			}
			return strings.Join(ss, "")
		}
	}(),
	"replace": func(s, old, new string) string {
		return strings.ReplaceAll(s, old, new)
	},
}
