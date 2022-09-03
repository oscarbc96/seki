package cloudformation

import (
	"github.com/awslabs/goformation/v6"
	"github.com/oscarbc96/seki/pkg/check"
	"github.com/oscarbc96/seki/pkg/load"
)

type S3BucketHasAllBlockPublicAccessEnabled struct{}

func (S3BucketHasAllBlockPublicAccessEnabled) Id() string { return "SK_1" }

func (S3BucketHasAllBlockPublicAccessEnabled) Name() string {
	return "Ensure S3 bucket has all \"block public access\" options enabled"
}

func (S3BucketHasAllBlockPublicAccessEnabled) Description() string { return "Description" }

func (S3BucketHasAllBlockPublicAccessEnabled) Severity() check.Severity { return check.Medium }

func (S3BucketHasAllBlockPublicAccessEnabled) Controls() map[string][]string {
	return map[string][]string{}
}

func (S3BucketHasAllBlockPublicAccessEnabled) Tags() []string { return []string{} }

func (S3BucketHasAllBlockPublicAccessEnabled) RemediationDoc() string {
	return "https://sekisecurity.com/docs/"
}

func (S3BucketHasAllBlockPublicAccessEnabled) InputTypes() []load.DetectedType {
	return []load.DetectedType{load.DetectedAwsCloudformation}
}

func (c S3BucketHasAllBlockPublicAccessEnabled) Run(f load.Input) check.CheckResult {
	template, err := goformation.Open(f.Path())
	if err != nil {
		return check.NewSkipCheckResultWithError(c, err)
	}

	var locations []load.Range
	for _, resource := range template.GetAllS3BucketResources() {
		if resource.PublicAccessBlockConfiguration == nil {
			locations = append(locations, load.Range{}) // TODO implement range
		} else {
			pabc := resource.PublicAccessBlockConfiguration
			if (pabc.BlockPublicAcls == nil || !*pabc.BlockPublicAcls) ||
				(pabc.BlockPublicPolicy == nil || !*pabc.BlockPublicPolicy) ||
				(pabc.IgnorePublicAcls == nil || !*pabc.IgnorePublicAcls) ||
				(pabc.RestrictPublicBuckets == nil || !*pabc.RestrictPublicBuckets) {
				locations = append(locations, load.Range{}) // TODO implement range
			}
		}
	}

	if len(locations) != 0 {
		return check.NewFailCheckResult(c, locations)
	}

	return check.NewPassCheckResult(c)
}
