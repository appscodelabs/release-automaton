# Release {{ .Release }}

{{ range $p := .Projects }}
## [{{ trimPrefix $p.URL "github.com/" }}](https://{{ $p.URL }})
{{ range $r := $p.Releases }}
### [{{ $r.Tag }}](https://{{ $p.URL }}/releases/tag/{{ $r.Tag }})

{{ range $c := $r.Commits -}}
 - [{{ slice $c.SHA 0 8 }}](https://{{ $p.URL }}/commit/{{ $c.SHA }}) {{ $c.Subject }}
{{ end }}
{{ end }}
{{ end }}
