package load

import "encoding/json"

type Position struct {
	Line   int
	Column int
}

func (p *Position) IsEmpty() bool {
	return p.Line <= 0 && p.Column <= 0
}

func (p *Position) MarshalJSON() ([]byte, error) {
	if p.Column <= 0 {
		return json.Marshal(&struct {
			Line int `json:"line"`
		}{
			Line: p.Line,
		})
	}

	return json.Marshal(&struct {
		Line   int `json:"line"`
		Column int `json:"column"`
	}{
		Line:   p.Line,
		Column: p.Column,
	})
}
