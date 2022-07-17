package check

import (
	"encoding/json"
)

type ErrorWrapper struct {
	Err error
}

func (ew *ErrorWrapper) IsEmpty() bool {
	return ew.Err == nil
}

func (ew *ErrorWrapper) MarshalJSON() ([]byte, error) {
	return json.Marshal(ew.Err.Error())
}
