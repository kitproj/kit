package types

import (
	"bytes"
	"encoding/csv"
	"encoding/json"
	"strings"
)

type Strings []string

func (p *Strings) UnmarshalJSON(data []byte) error {
	if data[0] == '[' {
		var x []string
		if err := json.Unmarshal(data, &x); err != nil {
			return err
		}
		for _, s := range x {
			*p = append(*p, s)
		}
		return nil
	}
	r := csv.NewReader(bytes.NewBuffer(data))
	r.Comma = ' '
	x, err := r.Read()
	for _, s := range x {
		*p = append(*p, s)
	}
	return err
}

func (p Strings) MarshalJSON() ([]byte, error) {
	b := new(bytes.Buffer)
	r := csv.NewWriter(b)
	r.Comma = ' '
	if err := r.WriteAll([][]string{p}); err != nil {
		return nil, err
	}
	return json.Marshal(strings.TrimSpace(b.String()))
}
