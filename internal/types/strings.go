package types

import (
	"bytes"
	"encoding/csv"
	"encoding/json"
	"fmt"
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
	var s string
	if err := json.Unmarshal(data, &s); err != nil {
		return err
	}
	r := csv.NewReader(bytes.NewBufferString(s))
	r.Comma = ' '
	x, err := r.Read()
	if err != nil {
		return fmt.Errorf("failed to read string (use double-quotes, not single-quotes): %w", err)
	}
	for _, s := range x {
		*p = append(*p, s)
	}
	return nil
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
