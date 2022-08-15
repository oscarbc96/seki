module github.com/oscarbc96/seki

go 1.18

require (
	github.com/awslabs/goformation/v6 v6.6.5
	github.com/distribution/distribution v2.8.1+incompatible
	github.com/jedib0t/go-pretty/v6 v6.3.6
	github.com/moby/buildkit v0.10.3
	github.com/rs/zerolog v1.27.0
	github.com/samber/lo v1.27.0
	github.com/spf13/afero v1.9.2
	github.com/spf13/cobra v1.5.0
	github.com/stretchr/testify v1.7.4
)

require (
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/mattn/go-runewidth v0.0.13 // indirect
	github.com/pmezard/go-difflib v1.0.0 // indirect
	github.com/rivo/uniseg v0.3.4 // indirect
	golang.org/x/exp v0.0.0-20220722155223-a9213eeb770e // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)

require (
	github.com/agext/levenshtein v1.2.3 // indirect
	github.com/containerd/typeurl v1.0.2 // indirect
	github.com/docker/distribution v2.8.1+incompatible // indirect
	github.com/docker/docker v20.10.17+incompatible // indirect; master (v22.xx-dev), see replace()
	github.com/docker/go-connections v0.4.0 // indirect
	github.com/docker/go-units v0.4.0 // indirect
	github.com/gogo/protobuf v1.3.2 // indirect
	github.com/golang/protobuf v1.5.2 // indirect
	github.com/inconshreveable/mousetrap v1.0.1 // indirect
	github.com/mattn/go-colorable v0.1.12 // indirect
	github.com/mattn/go-isatty v0.0.14 // indirect
	github.com/opencontainers/go-digest v1.0.0 // indirect
	github.com/opencontainers/image-spec v1.0.3-0.20211202183452-c5a74bcca799 // indirect
	github.com/pkg/errors v0.9.1
	github.com/sanathkr/go-yaml v0.0.0-20170819195128-ed9d249f429b // indirect
	github.com/sanathkr/yaml v0.0.0-20170819201035-0056894fa522 // indirect
	github.com/spf13/pflag v1.0.5 // indirect
	golang.org/x/sys v0.0.0-20220811171246-fbc7d0a398ab // indirect
	golang.org/x/text v0.3.7 // indirect
	google.golang.org/protobuf v1.28.1 // indirect
	gotest.tools/v3 v3.3.0 // indirect
)

replace github.com/docker/docker => github.com/docker/docker v20.10.3-0.20220414164044-61404de7df1a+incompatible
