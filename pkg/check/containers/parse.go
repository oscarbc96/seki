package containers

import (
	"github.com/distribution/distribution/v3/reference"
	"github.com/moby/buildkit/frontend/dockerfile/instructions"
	"github.com/moby/buildkit/frontend/dockerfile/parser"
	"github.com/oscarbc96/seki/pkg/load"
	"github.com/pkg/errors"
	"github.com/samber/lo"
)

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
