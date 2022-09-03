package report

import (
	"encoding/json"
	"github.com/pkg/errors"
	"strings"
)

type Formatter func(inputReports []InputReport) (string, error)

type Format int

const (
	JSON Format = iota
	Console
)

var formatToString = map[Format]string{
	JSON:    "json",
	Console: "console",
}

var DefaultFormat = formatToString[Console]

func (f Format) String() string {
	return formatToString[f]
}

func (f Format) MarshalJSON() ([]byte, error) {
	return json.Marshal(f.String())
}

func (f Format) GetFormatter() Formatter {
	switch f {
	case JSON:
		return JSONFormatter
	case Console:
		return ConsoleFormatter
	}
	return nil
}

func FormatterFromString(name string) (Formatter, error) {
	lower := strings.ToLower(name)
	for format, id := range formatToString {
		if id == lower {
			return format.GetFormatter(), nil
		}
	}
	return nil, errors.New("Unknown Format")
}
