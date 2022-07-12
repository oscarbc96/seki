package check

import "encoding/json"

type Status int

const (
	PASS Status = iota
	FAIL
	SKIP
)

var statusToString = map[Status]string{
	PASS: "PASS",
	FAIL: "FAIL",
	SKIP: "SKIP",
}

func (s Status) String() string {
	return statusToString[s]
}

func (s Status) MarshalJSON() ([]byte, error) {
	return json.Marshal(s.String())
}
