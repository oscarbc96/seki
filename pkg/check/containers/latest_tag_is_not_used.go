package containers

import (
	"github.com/oscarbc96/seki/pkg/check"
	"github.com/oscarbc96/seki/pkg/load"
)

type LatestTagIsNotUsed struct{}

func (LatestTagIsNotUsed) Id() string { return "SK_4" }

func (LatestTagIsNotUsed) Name() string {
	return "Ensure the base image uses a non latest version tag"
}

func (LatestTagIsNotUsed) Description() string { return "Description" }

func (LatestTagIsNotUsed) Severity() check.Severity { return check.Medium }

func (LatestTagIsNotUsed) Controls() map[string][]string {
	return map[string][]string{}
}

func (LatestTagIsNotUsed) Tags() []string { return []string{"docker"} }

func (LatestTagIsNotUsed) RemediationDoc() string { return "https://sekisecurity.com/docs/" }

func (LatestTagIsNotUsed) InputTypes() []load.DetectedType {
	return []load.DetectedType{load.DetectedContainerDockerfile}
}

func (c LatestTagIsNotUsed) Run(f load.Input) check.CheckResult {
	stages, _, err := parseDockerInstructions(f)
	if err != nil {
		return check.NewSkipCheckResultWithError(c, err)
	}

	dockerImages, err := parseDockerImagesFromStages(stages)
	if err != nil {
		return check.NewSkipCheckResultWithError(c, err)
	}

	var locations []load.Range
	for _, dockerImage := range dockerImages {
		if dockerImage.Tag == "latest" {
			locations = append(locations, dockerImage.Location)
		}
	}
	if len(locations) != 0 {
		return check.NewFailCheckResult(c, locations)
	}

	return check.NewPassCheckResult(c)
}
