package slomad

import (
	nomadStructs "github.com/hashicorp/nomad/nomad/structs"
)

type Port struct {
	Label  string
	To     int
	From   int
	Static bool
}

type PortParams struct {
	Label string
	To    int
	From  int
}

func NewPort(p PortParams) *Port {
	return &Port{
		Label:  p.Label,
		To:     p.To,
		From:   p.From,
		Static: p.To == p.From,
	}
}

func ToNomadPort(p *Port) nomadStructs.Port {
	return nomadStructs.Port{
		Label: p.Label,
		To:    p.To,
		Value: p.From,
	}
}

func NewPorts(args []PortParams) []*Port {
	ports := []*Port{}
	for _, arg := range args {
		ports = append(ports, NewPort(arg))
	}
	return ports
}

func ToPortMap(ports []*Port) map[string][]*Port {
	stat := []*Port{}
	dynm := []*Port{}

	for _, p := range ports {
		if p.Static {
			stat = append(stat, p)
		} else {
			dynm = append(dynm, p)
		}
	}

	portMap := map[string][]*Port{
		"static":  stat,
		"dynamic": dynm,
	}

	return portMap
}

func ToNomadPortMap(ports []*Port) map[string][]nomadStructs.Port {
	stat := []nomadStructs.Port{}
	dynm := []nomadStructs.Port{}

	for _, p := range ports {
		np := ToNomadPort(p)
		if p.Static {
			stat = append(stat, np)
		} else {
			dynm = append(dynm, np)
		}
	}

	portMap := map[string][]nomadStructs.Port{
		"static":  stat,
		"dynamic": dynm,
	}

	return portMap
}

func ExtractLabels(ports []*Port) []string {
	labels := []string{}
	for _, p := range ports {
		labels = append(labels, p.Label)
	}
	return labels
}
