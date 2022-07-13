package check

import (
	"github.com/samber/lo"
	"testing"
)

func TestCheckIdsAreUnique(t *testing.T) {
	ids := lo.Map[Check, string](allChecks, func(c Check, _ int) string {
		return c.Id()
	})
	duplicatedIds := lo.FindDuplicates[string](ids)

	if len(duplicatedIds) != 0 {
		t.Errorf("Found %v duplicated", duplicatedIds)
	}
}
