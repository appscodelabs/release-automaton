module github.com/appscodelabs/release-automaton

go 1.14

require (
	github.com/Masterminds/goutils v1.1.0 // indirect
	github.com/Masterminds/semver v1.5.0
	github.com/Masterminds/sprig v2.22.0+incompatible
	github.com/codeskyblue/go-sh v0.0.0-00010101000000-000000000000
	github.com/google/go-github/v32 v32.0.0
	github.com/google/uuid v1.1.1
	github.com/hashicorp/go-getter v1.4.1
	github.com/huandu/xstrings v1.3.2 // indirect
	github.com/imdario/mergo v0.3.9 // indirect
	github.com/keighl/metabolize v0.0.0-20150915210303-97ab655d4034
	github.com/mitchellh/copystructure v1.0.0 // indirect
	github.com/spf13/pflag v1.0.5
	github.com/tamalsaha/go-oneliners v0.0.0-20190126213733-6d24eabef827
	golang.org/x/mod v0.3.0
	golang.org/x/oauth2 v0.0.0-20190604053449-0f29369cfe45
	gomodules.xyz/envsubst v0.1.0
	k8s.io/apimachinery v0.18.3
	sigs.k8s.io/yaml v1.2.0
)

replace github.com/codeskyblue/go-sh => github.com/gomodules/go-sh v0.0.0-20200616225555-bfeba62378c9
