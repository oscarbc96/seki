package cloudformation

import (
	"github.com/awslabs/goformation/v6"
	"github.com/oscarbc96/seki/pkg/check"
	"github.com/oscarbc96/seki/pkg/load"
)

type S3BucketVersioningIsEnabled struct{}

func (S3BucketVersioningIsEnabled) Id() string { return "SK_2" }

func (S3BucketVersioningIsEnabled) Name() string {
	return "Ensure S3 bucket versioning is enabled"
}

func (S3BucketVersioningIsEnabled) Description() string {
	return "S3 bucket versioning is used to protect data availability and integrity. By enabling object versioning, data is protected from overwrites and deletions."
}

func (S3BucketVersioningIsEnabled) Severity() check.Severity { return check.Medium }

func (S3BucketVersioningIsEnabled) Controls() map[string][]string {
	return map[string][]string{}
}

func (S3BucketVersioningIsEnabled) Tags() []string { return []string{} }

func (S3BucketVersioningIsEnabled) RemediationDoc() string {
	return "https://sekisecurity.com/docs/"
}

func (S3BucketVersioningIsEnabled) InputTypes() []load.DetectedType {
	return []load.DetectedType{load.DetectedAwsCloudformation}
}

func (c S3BucketVersioningIsEnabled) Run(f load.Input) check.CheckResult {
	template, err := goformation.Open(f.Path())
	if err != nil {
		return check.NewSkipCheckResultWithError(c, err)
	}

	var locations []load.Range
	for _, resource := range template.GetAllS3BucketResources() {
		if resource.VersioningConfiguration != nil && resource.VersioningConfiguration.Status != "Enabled" {
			locations = append(locations, load.Range{}) // TODO implement range
		}
	}

	if len(locations) != 0 {
		return check.NewFailCheckResult(c, locations)
	}

	return check.NewPassCheckResult(c)
}
