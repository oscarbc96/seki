package github_actions

import "github.com/oscarbc96/seki/pkg/check"

var Checks []check.Check

func init() {
	Checks = []check.Check{
		new(Patata),
	}
}
