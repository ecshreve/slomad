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

func NewPort(l string, t, f int) *Port {
	return &Port{
		Label:  l,
		To:     t,
		From:   f,
		Static: t == f && t != 0,
	}
}

func ToNomadPort(p *Port) nomadStructs.Port {
	return nomadStructs.Port{
		Label: p.Label,
		To:    p.To,
		Value: p.From,
	}
}

func basicPort(val int) *Port {
	return &Port{
		Label:  "http",
		To:     val,
		From:   0,
		Static: false,
	}
}

func BasicPortConfig(val int) []*Port {
	return []*Port{basicPort(val)}
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

func getNetworks(ports []*Port) []*nomadStructs.NetworkResource {
	portMap := ToNomadPortMap(ports)
	return []*nomadStructs.NetworkResource{
		{
			ReservedPorts: portMap["static"],
			DynamicPorts:  portMap["dynamic"],
		},
	}
}
