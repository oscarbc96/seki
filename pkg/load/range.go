package load

import "encoding/json"

type Range struct {
	Start Position
	End   Position
}

func (r *Range) IsEmpty() bool {
	return r.Start.IsEmpty() && r.End.IsEmpty()
}

func (r *Range) MarshalJSON() ([]byte, error) {
	if r.Start.IsEmpty() && !r.End.IsEmpty() {
		return json.Marshal(&struct {
			End Position `json:"end"`
		}{
			End: r.End,
		})
	}

	if !r.Start.IsEmpty() && r.End.IsEmpty() {
		return json.Marshal(&struct {
			Start Position `json:"start"`
		}{
			Start: r.Start,
		})
	}

	return json.Marshal(&struct {
		Start Position `json:"start"`
		End   Position `json:"end"`
	}{
		Start: r.Start,
		End:   r.End,
	})
}
