package check

import (
	"github.com/distribution/distribution/reference"
	"github.com/moby/buildkit/frontend/dockerfile/instructions"
	"github.com/moby/buildkit/frontend/dockerfile/parser"
	"github.com/oscarbc96/seki/pkg/load"
	"github.com/oscarbc96/seki/utils"
	"github.com/pkg/errors"
	"github.com/samber/lo"
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

type DockerImage struct {
	Digest   string
	Image    string
	Platform string
	Registry string
	Tag      string
	Location load.Range
}

func parseDockerInstructions(f load.Input) (stages []instructions.Stage, metaArgs []instructions.ArgCommand, err error) {
	reader, err := f.Open()
	if err != nil {
		return nil, nil, err
	}
	defer reader.Close()

	parsedDockerfile, err := parser.Parse(reader)
	if err != nil {
		return nil, nil, err
	}

	return instructions.Parse(parsedDockerfile.AST)
}

func parseDockerImagesFromStages(stages []instructions.Stage) ([]DockerImage, error) {
	stageNames := lo.Map[instructions.Stage, string](stages, func(stage instructions.Stage, _ int) string {
		return stage.Name
	})

	var dockerImages []DockerImage
	for _, stage := range stages {
		// Ignoring if stage inherits from a previous stage
		if lo.Contains[string](stageNames, stage.BaseName) {
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

		cmdLocation := stage.Location[0] // TODO validate the hardcoded 0
		dockerImages = append(dockerImages, DockerImage{
			Digest:   digest,
			Image:    reference.Path(ref),
			Platform: stage.Platform,
			Registry: reference.Domain(ref),
			Tag:      ref.(reference.Tagged).Tag(),
			Location: load.Range{
				Start: load.Position{
					Line:   cmdLocation.Start.Line,
					Column: cmdLocation.Start.Character,
				},
				End: load.Position{
					Line:   cmdLocation.End.Line,
					Column: cmdLocation.End.Character,
				},
			},
		})
	}
	return dockerImages, nil
}

type ContainersDockerfileDockerHubRateLimit struct{}

func (ContainersDockerfileDockerHubRateLimit) Id() string { return "SK_3" }

func (ContainersDockerfileDockerHubRateLimit) Name() string { return "Name" }

func (ContainersDockerfileDockerHubRateLimit) Description() string { return "Description" }

func (ContainersDockerfileDockerHubRateLimit) Severity() Severity { return Informational }

func (ContainersDockerfileDockerHubRateLimit) Controls() map[string][]string {
	return map[string][]string{}
}

func (ContainersDockerfileDockerHubRateLimit) Tags() []string { return []string{"docker"} }

func (ContainersDockerfileDockerHubRateLimit) RemediationDoc() string { return "RemediationDoc" }

func (ContainersDockerfileDockerHubRateLimit) InputTypes() []load.DetectedType {
	return []load.DetectedType{load.DetectedContainerDockerfile}
}

func (c ContainersDockerfileDockerHubRateLimit) Run(f load.Input) CheckResult {
	stages, _, err := parseDockerInstructions(f)
	if err != nil {
		return NewSkipCheckResultWithError(c, err)
	}

	dockerImages, err := parseDockerImagesFromStages(stages)
	if err != nil {
		return NewSkipCheckResultWithError(c, err)
	}

	var locations []load.Range
	for _, dockerImage := range dockerImages {
		if dockerImage.Registry == "docker.io" {
			locations = append(locations, dockerImage.Location)
		}
	}
	if len(locations) != 0 {
		return NewFailCheckResult(c, locations)
	}

	return NewPassCheckResult(c)
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

func (c ContainersDockerfileLatestTag) Run(f load.Input) CheckResult {
	stages, _, err := parseDockerInstructions(f)
	if err != nil {
		return NewSkipCheckResultWithError(c, err)
	}

	dockerImages, err := parseDockerImagesFromStages(stages)
	if err != nil {
		return NewSkipCheckResultWithError(c, err)
	}

	var locations []load.Range
	for _, dockerImage := range dockerImages {
		if dockerImage.Tag == "latest" {
			locations = append(locations, dockerImage.Location)
		}
	}
	if len(locations) != 0 {
		return NewFailCheckResult(c, locations)
	}

	return NewPassCheckResult(c)
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

func (c ContainersDockerfileAddExists) Run(f load.Input) CheckResult {
	stages, _, err := parseDockerInstructions(f)
	if err != nil {
		return NewSkipCheckResultWithError(c, err)
	}

	var locations []load.Range
	for _, stage := range stages {
		for _, command := range stage.Commands {
			if _, isAddCommand := command.(*instructions.AddCommand); isAddCommand {
				cmdLocation := command.Location()[0] // TODO validate the hardcoded 0
				locations = append(locations, load.Range{
					Start: load.Position{
						Line:   cmdLocation.Start.Line,
						Column: 1,
					},
					End: load.Position{
						Line:   cmdLocation.End.Line,
						Column: len("ADD"),
					},
				})
			}
		}
	}

	if len(locations) != 0 {
		return NewFailCheckResult(c, locations)
	}

	return NewPassCheckResult(c)
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

func (c ContainersDockerfileRootUser) Run(f load.Input) CheckResult {
	stages, _, err := parseDockerInstructions(f)
	if err != nil {
		return NewSkipCheckResultWithError(c, err)
	}

	var locations []load.Range
	for _, stage := range stages {
		for _, command := range stage.Commands {
			if command, isUserCommand := command.(*instructions.UserCommand); isUserCommand {
				if command.User == "root" {
					colStart, colEnd := utils.FindStartAndEndColumn(command.String(), "root")
					cmdLocation := command.Location()[0] // TODO validate the hardcoded 0
					locations = append(locations, load.Range{
						Start: load.Position{
							Line:   cmdLocation.Start.Line,
							Column: colStart,
						},
						End: load.Position{
							Line:   cmdLocation.End.Line,
							Column: colEnd,
						},
					})
				}
			}
		}
	}

	if len(locations) != 0 {
		return NewFailCheckResult(c, locations)
	}

	return NewPassCheckResult(c)
}
