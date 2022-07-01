package check

import (
	check_aws_cloudformation_s3 "github.com/oscarbc96/seki/pkg/check/aws/cloudformation/s3"
	check_docker_dockerfile "github.com/oscarbc96/seki/pkg/check/docker/dockerfile"
	"github.com/oscarbc96/seki/pkg/result"
)

type Check func() ([]result.CheckResult, error)

var Checkers = map[string]Check{
	"CheckS3BucketPublicReadAcl":  check_aws_cloudformation_s3.CheckS3BucketPublicReadAcl,
	"CheckS3ObjectVersioningRule": check_aws_cloudformation_s3.CheckS3ObjectVersioningRule,
	"CheckDockerHubRateLimit":     check_docker_dockerfile.CheckDockerHubRateLimit,
	"CheckLatestTag":              check_docker_dockerfile.CheckLatestTag,
}
