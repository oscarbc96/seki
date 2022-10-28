package github_actions

import (
	"github.com/oscarbc96/seki/pkg/check"
	"github.com/oscarbc96/seki/pkg/load"
)

type Patata struct{}

func (Patata) Id() string { return "SK_1" }

func (Patata) Name() string {
	return "Ensure S3 bucket has all \"block public access\" options enabled"
}

func (Patata) Description() string { return "Description" }

func (Patata) Severity() check.Severity { return check.Medium }

func (Patata) Controls() map[string][]string {
	return map[string][]string{}
}

func (Patata) Tags() []string { return []string{} }

func (Patata) RemediationDoc() string {
	return "https://sekisecurity.com/docs/"
}

func (Patata) InputTypes() []load.DetectedType {
	return []load.DetectedType{load.DetectedCICDGitHubActions}
}

func (c Patata) Run(f load.Input) check.CheckResult {

	return check.NewPassCheckResult(c)
}
