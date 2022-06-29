package check_aws_cloudformation_s3

import (
	"fmt"
	"github.com/awslabs/goformation/v6"
	"github.com/oscarbc96/seki/pkg/check"
	"github.com/oscarbc96/seki/pkg/result"
)

func init() {
	check.Analysers["CheckS3BucketPublicReadAcl"] = CheckS3BucketPublicReadAcl
	check.Analysers["CheckS3ObjectVersioningRule"] = CheckS3ObjectVersioningRule
}

func CheckS3BucketPublicReadAcl() (*result.CheckResult, error) {
	template, err := goformation.Open("test.json")
	if err != nil {
		return nil, err
	}

	for logicalID, resource := range template.GetAllS3BucketResources() {
		if resource.AccessControl != nil && *resource.AccessControl == "PublicRead" {
			return &result.CheckResult{Result: result.FAIL, Severity: result.Medium, Message: fmt.Sprintf("S3 Bucket %s should not have a public-read acl", logicalID)}, nil
		}
	}

	return &result.CheckResult{Result: result.PASS, Severity: result.Medium}, nil
}

func CheckS3ObjectVersioningRule() (*result.CheckResult, error) {
	template, err := goformation.Open("test.json")
	if err != nil {
		return nil, err
	}

	for logicalID, resource := range template.GetAllS3BucketResources() {
		if resource.VersioningConfiguration != nil && resource.VersioningConfiguration.Status != "Enabled" {
			// TODO Chec
			return &result.CheckResult{Result: result.FAIL, Severity: result.Medium, Message: fmt.Sprintf("S3 Bucket %s is required to have object versioning enabled", logicalID)}, nil
		}
	}

	return &result.CheckResult{Result: result.PASS, Severity: result.Medium}, nil
}
