package containers

import "github.com/oscarbc96/seki/pkg/check"

var AllChecks []check.Check

func init() {
	AllChecks = append(
		AllChecks,
		new(AdviseDockerHubRateLimit),
		new(LatestTagIsNotUsed),
		new(PreferCopyOverAdd),
		new(LastUserIsNotRoot),
	)
}
