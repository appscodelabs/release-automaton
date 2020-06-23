# {{ .ProductLine }} Releases

| {{ .ProductLine }} Version | Release Date | User Guide | Changelog | Kubernetes Version |
|--------------------------- | ------------ | ---------- | --------- | ------------------ |
{{ range $r :=  .Releases -}}
| [{{ $r.Release }}](https://{{ $.URL }}/tag/{{ $r.Release }}) | {{ $r.ReleaseDate }} | [User Guide]({ $r.ChangelogLink }}) | [CHANGELOG]({{ $r.UserGuideLink }}) | {{ $r.KubernetesVersion }} |
{{ end }}
