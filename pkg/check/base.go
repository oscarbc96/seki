package check

import (
	check_aws_cloudformation_s3 "github.com/oscarbc96/seki/pkg/check/aws/cloudformation/s3"
	check_docker_dockerfile "github.com/oscarbc96/seki/pkg/check/docker/dockerfile"
	"github.com/oscarbc96/seki/pkg/load"
	"github.com/oscarbc96/seki/pkg/result"
)

type Check struct {
	Id          string
	Title       string
	Description string
	Severity    string
	Controls    map[string][]string

	//Providers   []string
	//Categories  string
	Run CheckFunction
}

type CheckFunction func(file load.InputFile) ([]result.CheckResult, error)

var Checkers = map[string][]CheckFunction{
	"cloudformation": {
		check_aws_cloudformation_s3.CheckS3BucketPublicReadAcl,
		check_aws_cloudformation_s3.CheckS3ObjectVersioningRule,
	},
	"dockerfile": {
		check_docker_dockerfile.CheckDockerHubRateLimit,
		check_docker_dockerfile.CheckLatestTag,
	},
}
