package check

import (
	"github.com/awslabs/goformation/v6"
	"github.com/oscarbc96/seki/pkg/load"
)

func init() {
	allChecks = append(
		allChecks,
		new(AWSCloudformationS3BucketPublicReadAcl),
		new(AWSCloudformationS3ObjectVersioningRule),
	)
}

type AWSCloudformationS3BucketPublicReadAcl struct{}

func (AWSCloudformationS3BucketPublicReadAcl) Id() string { return "SK_1" }

func (AWSCloudformationS3BucketPublicReadAcl) Name() string { return "Name" }

func (AWSCloudformationS3BucketPublicReadAcl) Description() string { return "Description" }

func (AWSCloudformationS3BucketPublicReadAcl) Severity() Severity { return Medium }

func (AWSCloudformationS3BucketPublicReadAcl) Controls() map[string][]string {
	return map[string][]string{}
}

func (AWSCloudformationS3BucketPublicReadAcl) Tags() []string { return []string{} }

func (AWSCloudformationS3BucketPublicReadAcl) RemediationDoc() string { return "RemediationDoc" }

func (AWSCloudformationS3BucketPublicReadAcl) InputTypes() []load.DetectedType {
	return []load.DetectedType{load.DetectedAwsCloudformation}
}

func (c AWSCloudformationS3BucketPublicReadAcl) Run(f load.Input) (CheckResult, error) {
	template, err := goformation.Open(f.Path())
	if err != nil {
		return NewSkipCheckResult(c), err
	}

	var locations []load.Range
	for _, resource := range template.GetAllS3BucketResources() {
		if resource.AccessControl != nil && *resource.AccessControl == "PublicRead" {
			locations = append(locations, load.Range{}) // TODO implement range
		}
	}

	if len(locations) != 0 {
		return NewFailCheckResult(c, locations), nil
	}

	return NewPassCheckResult(c), nil
}

type AWSCloudformationS3ObjectVersioningRule struct{}

func (AWSCloudformationS3ObjectVersioningRule) Id() string { return "SK_2" }

func (AWSCloudformationS3ObjectVersioningRule) Name() string { return "Name" }

func (AWSCloudformationS3ObjectVersioningRule) Description() string { return "Description" }

func (AWSCloudformationS3ObjectVersioningRule) Severity() Severity { return Medium }

func (AWSCloudformationS3ObjectVersioningRule) Controls() map[string][]string {
	return map[string][]string{}
}

func (AWSCloudformationS3ObjectVersioningRule) Tags() []string { return []string{} }

func (AWSCloudformationS3ObjectVersioningRule) RemediationDoc() string { return "RemediationDoc" }

func (AWSCloudformationS3ObjectVersioningRule) InputTypes() []load.DetectedType {
	return []load.DetectedType{load.DetectedAwsCloudformation}
}

func (c AWSCloudformationS3ObjectVersioningRule) Run(f load.Input) (CheckResult, error) {
	template, err := goformation.Open(f.Path())
	if err != nil {
		return NewSkipCheckResult(c), err
	}

	var locations []load.Range
	for _, resource := range template.GetAllS3BucketResources() {
		if resource.VersioningConfiguration != nil && resource.VersioningConfiguration.Status != "Enabled" {
			locations = append(locations, load.Range{}) // TODO implement range
		}
	}

	if len(locations) != 0 {
		return NewFailCheckResult(c, locations), nil
	}

	return NewPassCheckResult(c), nil
}
