package check

import (
	"github.com/awslabs/goformation/v6"
	"github.com/oscarbc96/seki/pkg/load"
)

func init() {
	allChecks = append(
		allChecks,
		new(CheckAWSCloudformationS3BucketPublicReadAcl),
		new(CheckAWSCloudformationS3ObjectVersioningRule),
	)
}

type CheckAWSCloudformationS3BucketPublicReadAcl struct{}

func (CheckAWSCloudformationS3BucketPublicReadAcl) Id() string { return "SK_1" }

func (CheckAWSCloudformationS3BucketPublicReadAcl) Name() string { return "Name" }

func (CheckAWSCloudformationS3BucketPublicReadAcl) Description() string { return "Description" }

func (CheckAWSCloudformationS3BucketPublicReadAcl) Severity() Severity { return Medium }

func (CheckAWSCloudformationS3BucketPublicReadAcl) Controls() map[string][]string {
	return map[string][]string{}
}

func (CheckAWSCloudformationS3BucketPublicReadAcl) Tags() []string { return []string{} }

func (CheckAWSCloudformationS3BucketPublicReadAcl) RemediationDoc() string { return "RemediationDoc" }

func (CheckAWSCloudformationS3BucketPublicReadAcl) InputTypes() []load.DetectedType {
	return []load.DetectedType{load.DetectedAwsCloudformation}
}

func (CheckAWSCloudformationS3BucketPublicReadAcl) Run(f load.Input) (CheckResult, error) {
	template, err := goformation.Open(f.Path())
	if err != nil {
		return NewSkipCheckResult(), err
	}

	var locations []load.Range
	for _, resource := range template.GetAllS3BucketResources() {
		if resource.AccessControl != nil && *resource.AccessControl == "PublicRead" {
			locations = append(locations, load.Range{}) // TODO implement range
		}
	}

	if len(locations) != 0 {
		return NewFailCheckResult(locations), nil
	}

	return NewPassCheckResult(), nil
}

type CheckAWSCloudformationS3ObjectVersioningRule struct{}

func (CheckAWSCloudformationS3ObjectVersioningRule) Id() string { return "SK_2" }

func (CheckAWSCloudformationS3ObjectVersioningRule) Name() string { return "Name" }

func (CheckAWSCloudformationS3ObjectVersioningRule) Description() string { return "Description" }

func (CheckAWSCloudformationS3ObjectVersioningRule) Severity() Severity { return Medium }

func (CheckAWSCloudformationS3ObjectVersioningRule) Controls() map[string][]string {
	return map[string][]string{}
}

func (CheckAWSCloudformationS3ObjectVersioningRule) Tags() []string { return []string{} }

func (CheckAWSCloudformationS3ObjectVersioningRule) RemediationDoc() string { return "RemediationDoc" }

func (CheckAWSCloudformationS3ObjectVersioningRule) InputTypes() []load.DetectedType {
	return []load.DetectedType{load.DetectedAwsCloudformation}
}

func (CheckAWSCloudformationS3ObjectVersioningRule) Run(f load.Input) (CheckResult, error) {
	template, err := goformation.Open(f.Path())
	if err != nil {
		return NewSkipCheckResult(), err
	}

	var locations []load.Range
	for _, resource := range template.GetAllS3BucketResources() {
		if resource.VersioningConfiguration != nil && resource.VersioningConfiguration.Status != "Enabled" {
			locations = append(locations, load.Range{}) // TODO implement range
		}
	}

	if len(locations) != 0 {
		return NewFailCheckResult(locations), nil
	}

	return NewPassCheckResult(), nil
}
