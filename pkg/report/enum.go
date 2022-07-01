package report

import (
	"fmt"
	"github.com/oscarbc96/seki/pkg/result"
	"strings"
)

type Formater func(r []result.CheckResult) (string, error)

type Format int

const (
	JSON Format = iota
	SARIF
)

var FormatIDs = map[Format][]string{
	JSON:  {"json"},
	SARIF: {"sarif"},
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
	return -1, fmt.Errorf("Unknown Format: '%s'", name)
}

func GetFormater(f Format) (Formater, error) {
	switch f {
	case JSON:
		return JSONFormat, nil
	case SARIF:
		return SARIFReport, nil

	default:
		return nil, fmt.Errorf("Unrecognized formater: %s", FormatIDs[f])
	}
}
