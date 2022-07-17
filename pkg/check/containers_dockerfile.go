package check

import (
	"github.com/distribution/distribution/reference"
	"github.com/moby/buildkit/frontend/dockerfile/instructions"
	"github.com/moby/buildkit/frontend/dockerfile/parser"
	"github.com/oscarbc96/seki/pkg/load"
	"github.com/pkg/errors"
	"github.com/samber/lo"
	"io"
)

func init() {
	allChecks = append(
		allChecks,
		new(ContainersDockerfileDockerHubRateLimit),
		new(ContainersDockerfileLatestTag),
		new(ContainersDockerfileAddExists),
		new(ContainersDockerfileRootUser),
	)
}

type DockerStage struct {
	Digest   string
	Image    string
	Platform string
	Registry string
	Tag      string
	Location load.Range
}

func parseDockerStages(file io.Reader) ([]DockerStage, error) {
	parsedDockerfile, err := parser.Parse(file)
	if err != nil {
		return nil, err
	}

	stages, _, err := instructions.Parse(parsedDockerfile.AST)
	if err != nil {
		return nil, err
	}

	layerNames := lo.Map[instructions.Stage, string](stages, func(stage instructions.Stage, _ int) string {
		return stage.Name
	})

	var matches []DockerStage
	for _, stage := range stages {
		// Ignoring if stage inherits from a previous stage
		if lo.Contains[string](layerNames, stage.BaseName) {
			continue
		}

		ref, err := reference.ParseNormalizedNamed(stage.BaseName)
		if err != nil {
			return nil, errors.Wrapf(err, "Error parsing reference: %q is not a valid repository/tag", stage.BaseName)
		}
		ref = reference.TagNameOnly(ref)

		var digest string
		if canonicalReference, isCanonical := ref.(reference.Canonical); isCanonical {
			digest = canonicalReference.Digest().String()
		}

		matches = append(matches, DockerStage{
			Digest:   digest,
			Image:    reference.Path(ref),
			Platform: stage.Platform,
			Registry: reference.Domain(ref),
			Tag:      ref.(reference.Tagged).Tag(),
			Location: load.Range{
				Start: load.Position{
					Line:   stage.Location[0].Start.Line, // TODO validate the hardcoded 0
					Column: stage.Location[0].Start.Character,
				},
				End: load.Position{
					Line:   stage.Location[0].End.Line,
					Column: stage.Location[0].End.Character,
				},
			},
		})
	}
	return matches, nil
}

type ContainersDockerfileDockerHubRateLimit struct{}

func (ContainersDockerfileDockerHubRateLimit) Id() string { return "SK_3" }

func (ContainersDockerfileDockerHubRateLimit) Name() string { return "Name" }

func (ContainersDockerfileDockerHubRateLimit) Description() string { return "Description" }

func (ContainersDockerfileDockerHubRateLimit) Severity() Severity { return Medium }

func (ContainersDockerfileDockerHubRateLimit) Controls() map[string][]string {
	return map[string][]string{}
}

func (ContainersDockerfileDockerHubRateLimit) Tags() []string { return []string{} }

func (ContainersDockerfileDockerHubRateLimit) RemediationDoc() string { return "RemediationDoc" }

func (ContainersDockerfileDockerHubRateLimit) InputTypes() []load.DetectedType {
	return []load.DetectedType{load.DetectedContainerDockerfile}
}

func (ContainersDockerfileDockerHubRateLimit) Run(f load.Input) (CheckResult, error) {
	file, err := f.Open()
	if err != nil {
		return NewSkipCheckResult(), err
	}
	defer file.Close()

	dockerLayers, err := parseDockerStages(file)
	if err != nil {
		return NewSkipCheckResult(), err
	}

	var locations []load.Range
	for _, layer := range dockerLayers {
		if layer.Registry == "docker.io" {
			locations = append(locations, layer.Location)
		}
	}
	if len(locations) != 0 {
		return NewFailCheckResult(locations), nil
	}

	return NewPassCheckResult(), nil
}

type ContainersDockerfileLatestTag struct{}

func (ContainersDockerfileLatestTag) Id() string { return "SK_4" }

func (ContainersDockerfileLatestTag) Name() string {
	return "Ensure the base image uses a non latest version tag"
}

func (ContainersDockerfileLatestTag) Description() string { return "Description" }

func (ContainersDockerfileLatestTag) Severity() Severity { return Medium }

func (ContainersDockerfileLatestTag) Controls() map[string][]string {
	return map[string][]string{}
}

func (ContainersDockerfileLatestTag) Tags() []string { return []string{"docker"} }

func (ContainersDockerfileLatestTag) RemediationDoc() string { return "RemediationDoc" }

func (ContainersDockerfileLatestTag) InputTypes() []load.DetectedType {
	return []load.DetectedType{load.DetectedContainerDockerfile}
}

func (ContainersDockerfileLatestTag) Run(f load.Input) (CheckResult, error) {
	file, err := f.Open()
	if err != nil {
		return NewSkipCheckResult(), err
	}
	defer file.Close()

	dockerLayers, err := parseDockerStages(file)
	if err != nil {
		return NewSkipCheckResult(), err
	}

	var locations []load.Range
	for _, layer := range dockerLayers {
		if layer.Tag == "latest" {
			locations = append(locations, layer.Location)
		}
	}
	if len(locations) != 0 {
		return NewFailCheckResult(locations), nil
	}

	return NewPassCheckResult(), nil
}

type ContainersDockerfileAddExists struct{}

func (ContainersDockerfileAddExists) Id() string { return "SK_5" }

func (ContainersDockerfileAddExists) Name() string {
	return "Ensure that COPY is used instead of ADD"
}

func (ContainersDockerfileAddExists) Description() string { return "Description" }

func (ContainersDockerfileAddExists) Severity() Severity { return Medium }

func (ContainersDockerfileAddExists) Controls() map[string][]string {
	return map[string][]string{}
}

func (ContainersDockerfileAddExists) Tags() []string { return []string{"docker"} }

func (ContainersDockerfileAddExists) RemediationDoc() string { return "RemediationDoc" }

func (ContainersDockerfileAddExists) InputTypes() []load.DetectedType {
	return []load.DetectedType{load.DetectedContainerDockerfile}
}

func (ContainersDockerfileAddExists) Run(f load.Input) (CheckResult, error) {
	file, err := f.Open()
	if err != nil {
		return NewSkipCheckResult(), err
	}
	defer file.Close()

	parsedDockerfile, err := parser.Parse(file)
	if err != nil {
		return NewSkipCheckResult(), err
	}

	stages, _, err := instructions.Parse(parsedDockerfile.AST)
	if err != nil {
		return NewSkipCheckResult(), err
	}

	var locations []load.Range
	for _, stage := range stages {
		for _, command := range stage.Commands {
			if _, isAddCommand := command.(*instructions.AddCommand); isAddCommand {
				locations = append(locations, load.Range{
					Start: load.Position{
						Line:   command.Location()[0].Start.Line, // TODO validate the hardcoded 0
						Column: command.Location()[0].Start.Character,
					},
					End: load.Position{
						Line:   command.Location()[0].End.Line,
						Column: command.Location()[0].End.Character,
					},
				})
			}
		}
	}

	if len(locations) != 0 {
		return NewFailCheckResult(locations), nil
	}

	return NewPassCheckResult(), nil
}

type ContainersDockerfileRootUser struct{}

func (ContainersDockerfileRootUser) Id() string { return "SK_6" }

func (ContainersDockerfileRootUser) Name() string {
	return "Ensure the last USER is not root"
}

func (ContainersDockerfileRootUser) Description() string { return "Description" }

func (ContainersDockerfileRootUser) Severity() Severity { return Medium }

func (ContainersDockerfileRootUser) Controls() map[string][]string {
	return map[string][]string{}
}

func (ContainersDockerfileRootUser) Tags() []string { return []string{"docker"} }

func (ContainersDockerfileRootUser) RemediationDoc() string { return "RemediationDoc" }

func (ContainersDockerfileRootUser) InputTypes() []load.DetectedType {
	return []load.DetectedType{load.DetectedContainerDockerfile}
}

func (ContainersDockerfileRootUser) Run(f load.Input) (CheckResult, error) {
	file, err := f.Open()
	if err != nil {
		return NewSkipCheckResult(), err
	}
	defer file.Close()

	parsedDockerfile, err := parser.Parse(file)
	if err != nil {
		return NewSkipCheckResult(), err
	}

	stages, _, err := instructions.Parse(parsedDockerfile.AST)
	if err != nil {
		return NewSkipCheckResult(), err
	}

	var locations []load.Range
	for _, stage := range stages {
		for _, command := range stage.Commands {
			if command, isUserCommand := command.(*instructions.UserCommand); isUserCommand {
				if command.User == "root" {
					locations = append(locations, load.Range{
						Start: load.Position{
							Line:   command.Location()[0].Start.Line, // TODO validate the hardcoded 0
							Column: command.Location()[0].Start.Character,
						},
						End: load.Position{
							Line:   command.Location()[0].End.Line,
							Column: command.Location()[0].End.Character,
						},
					})
				}
			}
		}
	}

	if len(locations) != 0 {
		return NewFailCheckResult(locations), nil
	}

	return NewPassCheckResult(), nil
}
