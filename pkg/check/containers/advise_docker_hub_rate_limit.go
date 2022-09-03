package containers

import (
	"github.com/oscarbc96/seki/pkg/check"
	"github.com/oscarbc96/seki/pkg/load"
)

type AdviseDockerHubRateLimit struct{}

func (AdviseDockerHubRateLimit) Id() string { return "SK_3" }

func (AdviseDockerHubRateLimit) Name() string {
	return "Ensure the registry is not empty nor docker.io"
}

func (AdviseDockerHubRateLimit) Description() string { return "Description" }

func (AdviseDockerHubRateLimit) Severity() check.Severity { return check.Informational }

func (AdviseDockerHubRateLimit) Controls() map[string][]string {
	return map[string][]string{}
}

func (AdviseDockerHubRateLimit) Tags() []string { return []string{"docker"} }

func (AdviseDockerHubRateLimit) RemediationDoc() string {
	return "https://sekisecurity.com/docs/"
}

func (AdviseDockerHubRateLimit) InputTypes() []load.DetectedType {
	return []load.DetectedType{load.DetectedContainerDockerfile}
}

func (c AdviseDockerHubRateLimit) Run(f load.Input) check.CheckResult {
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
		if dockerImage.Registry == "docker.io" {
			locations = append(locations, dockerImage.Location)
		}
	}
	if len(locations) != 0 {
		return check.NewFailCheckResult(c, locations)
	}

	return check.NewPassCheckResult(c)
}
