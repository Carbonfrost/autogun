module github.com/Carbonfrost/autogun

go 1.26.4

require (
	github.com/Carbonfrost/joe-cli v0.9.1
	github.com/chromedp/cdproto v0.0.0-20250803210736-d308e07a266d
	github.com/chromedp/chromedp v0.14.2
	github.com/hashicorp/hcl/v2 v2.24.0
	github.com/onsi/ginkgo/v2 v2.32.0
	github.com/onsi/gomega v1.42.1
	github.com/spf13/afero v1.15.0
	github.com/zclconf/go-cty v1.17.0
	golang.org/x/term v0.44.0
)

require (
	codeberg.org/chavacava/garif v0.2.1 // indirect
	github.com/BurntSushi/toml v1.6.0 // indirect
	github.com/Masterminds/semver/v3 v3.5.0 // indirect
	github.com/agext/levenshtein v1.2.3 // indirect
	github.com/apparentlymart/go-textseg/v15 v15.0.0 // indirect
	github.com/chromedp/sysutil v1.1.0 // indirect
	github.com/cristalhq/acmd v0.12.0 // indirect
	github.com/fatih/color v1.19.0 // indirect
	github.com/fatih/structtag v1.2.0 // indirect
	github.com/go-critic/go-critic v0.12.0 // indirect
	github.com/go-json-experiment/json v0.0.0-20251027170946-4849db3c2f7e // indirect
	github.com/go-logr/logr v1.4.3 // indirect
	github.com/go-task/slim-sprig/v3 v3.0.0 // indirect
	github.com/go-toolsmith/astcast v1.1.0 // indirect
	github.com/go-toolsmith/astcopy v1.1.0 // indirect
	github.com/go-toolsmith/astequal v1.2.0 // indirect
	github.com/go-toolsmith/astfmt v1.1.0 // indirect
	github.com/go-toolsmith/astp v1.1.0 // indirect
	github.com/go-toolsmith/pkgload v1.2.2 // indirect
	github.com/go-toolsmith/strparse v1.1.0 // indirect
	github.com/go-toolsmith/typep v1.1.0 // indirect
	github.com/gobwas/httphead v0.1.0 // indirect
	github.com/gobwas/pool v0.2.1 // indirect
	github.com/gobwas/ws v1.4.0 // indirect
	github.com/google/go-cmp v0.7.0 // indirect
	github.com/google/pprof v0.0.0-20260604005048-7023385849c0 // indirect
	github.com/hashicorp/go-version v1.9.0 // indirect
	github.com/juju/ansiterm v1.0.0 // indirect
	github.com/kballard/go-shellquote v0.0.0-20180428030007-95032a82bc51 // indirect
	github.com/lunixbochs/vtclean v1.0.0 // indirect
	github.com/mattn/go-colorable v0.1.15 // indirect
	github.com/mattn/go-isatty v0.0.22 // indirect
	github.com/maxbrunsfeld/counterfeiter/v6 v6.12.2 // indirect
	github.com/mgechev/dots v1.0.0 // indirect
	github.com/mgechev/revive v1.15.0 // indirect
	github.com/mitchellh/go-ps v1.0.0 // indirect
	github.com/mitchellh/go-wordwrap v1.0.1 // indirect
	github.com/mitchellh/mapstructure v1.5.0 // indirect
	github.com/quasilyte/go-ruleguard v0.4.2 // indirect
	github.com/quasilyte/gogrep v0.5.0 // indirect
	github.com/quasilyte/regex/syntax v0.0.0-20210819130434-b3f0c404a727 // indirect
	github.com/quasilyte/stdinfo v0.0.0-20220114132959-f7386bf02567 // indirect
	go.yaml.in/yaml/v3 v3.0.4 // indirect
	golang.org/x/exp/typeparams v0.0.0-20260611194520-c48552f49976 // indirect
	golang.org/x/mod v0.37.0 // indirect
	golang.org/x/net v0.56.0 // indirect
	golang.org/x/sync v0.21.0 // indirect
	golang.org/x/sys v0.46.0 // indirect
	golang.org/x/text v0.38.0 // indirect
	golang.org/x/tools v0.47.0 // indirect
	honnef.co/go/tools v0.7.0 // indirect
)

tool (
	github.com/go-critic/go-critic/cmd/gocritic
	github.com/hashicorp/hcl/v2/cmd/hclfmt
	github.com/maxbrunsfeld/counterfeiter/v6
	github.com/mgechev/revive
	github.com/onsi/ginkgo/v2/ginkgo
	honnef.co/go/tools/cmd/staticcheck
)

ignore (
	.github
	bin
	eng
)
