package check_aws_cloudformation_s3

import (
	"fmt"
	"github.com/awslabs/goformation/v6"
	"github.com/oscarbc96/seki/pkg/load"
	"github.com/oscarbc96/seki/pkg/result"
)

func CheckS3BucketPublicReadAcl(f load.InputFile) ([]result.CheckResult, error) {
	template, err := goformation.Open("test.json")
	if err != nil {
		return nil, err
	}

	var output []result.CheckResult
	for logicalID, resource := range template.GetAllS3BucketResources() {
		if resource.AccessControl != nil && *resource.AccessControl == "PublicRead" {
			output = append(
				output,
				result.CheckResult{
					Result:   result.FAIL,
					Severity: result.Medium,
					Message:  fmt.Sprintf("S3 Bucket %s should not have a public-read acl", logicalID),
				},
			)
		}
	}

	return output, nil
}

func CheckS3ObjectVersioningRule(f load.InputFile) ([]result.CheckResult, error) {
	template, err := goformation.Open("test.json")
	if err != nil {
		return nil, err
	}

	var output []result.CheckResult
	for logicalID, resource := range template.GetAllS3BucketResources() {
		if resource.VersioningConfiguration != nil && resource.VersioningConfiguration.Status != "Enabled" {
			output = append(
				output,
				result.CheckResult{
					Result:   result.FAIL,
					Severity: result.Medium,
					Message:  fmt.Sprintf("S3 Bucket %s is required to have object versioning enabled", logicalID),
				},
			)
		}
	}

	return output, nil
}
