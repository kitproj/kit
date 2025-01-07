package types

import "encoding/json"

// A list of ports to expose.
type Ports []Port

func (p Ports) Len() int {
	return len(p)
}

func (p Ports) Less(i, j int) bool {
	return p[i].HostPort < p[j].HostPort
}

func (p Ports) Swap(i, j int) {
	p[i], p[j] = p[j], p[i]
}

func (p *Ports) UnmarshalJSON(data []byte) error {
	if data[0] == '[' {
		var x []Port
		if err := json.Unmarshal(data, &x); err != nil {
			return err
		}
		for _, port := range x {
			*p = append(*p, port)
		}
		return nil
	}
	var i int
	if err := json.Unmarshal(data, &i); err == nil {
		*p = append(*p, Port{ContainerPort: uint16(i)})
		return nil
	}
	var x = Strings{}
	if err := json.Unmarshal(data, &x); err != nil {
		return err
	}
	for _, port := range x {
		y := Port{}
		if err := y.Unstring(port); err != nil {
			return err
		}
		*p = append(*p, y)
	}

	return nil
}

func (p Ports) MarshalJSON() ([]byte, error) {
	var x Strings
	for _, port := range p {
		x = append(x, port.String())
	}
	return json.Marshal(x)
}

func (p Ports) Map() map[uint16]uint16 {
	m := map[uint16]uint16{}
	for _, port := range p {
		m[port.ContainerPort] = port.GetHostPort()
	}
	return m
}
