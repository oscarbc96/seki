package run

import (
	"github.com/oscarbc96/seki/pkg/check"
	"github.com/oscarbc96/seki/pkg/check/containers"
	"github.com/oscarbc96/seki/pkg/load"
	"github.com/samber/lo"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestCheckIdsAreUnique(t *testing.T) {
	ids := lo.Map[check.Check, string](AllChecks, func(c check.Check, _ int) string {
		return c.Id()
	})
	duplicatedIds := lo.FindDuplicates[string](ids)

	assert.Empty(t, duplicatedIds, "Found %v duplicated", duplicatedIds)
}

func TestNumberOfChecks(t *testing.T) {
	assert.Equal(t, 8, len(AllChecks))
}

func TestGetChecksFor(t *testing.T) {
	type args struct {
		tpe load.DetectedType
	}
	tests := []struct {
		name string
		args args
		want []check.Check
	}{
		{
			name: "Get docker checks",
			args: args{tpe: load.DetectedContainerDockerfile},
			want: []check.Check{
				new(containers.AdviseDockerHubRateLimit),
				new(containers.LatestTagIsNotUsed),
				new(containers.PreferCopyOverAdd),
				new(containers.LastUserIsNotRoot),
				new(containers.WorkdirPathMustBeAbsolute),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equalf(t, tt.want, GetChecksFor(tt.args.tpe), "GetChecksFor(%v)", tt.args.tpe)
		})
	}
}
