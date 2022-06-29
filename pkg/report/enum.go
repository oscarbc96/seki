package report

import (
	"fmt"
	"github.com/oscarbc96/seki/pkg/result"
	"strings"
)

type Formater func(r *result.CheckResult) (string, error)

type Format int

const (
	JSON Format = iota
)

var FormatIDs = map[Format][]string{
	JSON: {"json"},
}

var DefaultFormat = FormatIDs[JSON][0]

func FormatFromString(name string) (Format, error) {
	lower := strings.ToLower(name)
	for f, ids := range FormatIDs {
		for _, i := range ids {
			if lower == i {
				return f, nil
			}
		}
	}
	return -1, fmt.Errorf("Unknown Format String: '%s'", name)
}

func ValidateFormat(name string) error {
	if _, err := FormatFromString(name); err != nil {
		return err
	}
	return nil
}

func GetFormater(f Format) (Formater, error) {
	switch f {
	case JSON:
		return JSONFormat, nil
	default:
		return nil, fmt.Errorf("Unrecognized formater: %s", FormatIDs[f])
	}
}
