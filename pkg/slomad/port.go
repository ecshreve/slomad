package slomad

import (
	"fmt"

	nomadStructs "github.com/hashicorp/nomad/nomad/structs"
)

type Port struct {
	Label string
	To    int
	From  *int
}

func NewPort(label string, to int, from *int) *Port {
	p := &Port{
		Label: label,
		To:    to,
	}

	if from != nil {
		p.From = from
	}

	return p
}

func (p *Port) ToNomadPort() nomadStructs.Port {
	port := nomadStructs.Port{
		Label: p.Label,
		To:    p.To,
	}

	if p.From != nil {
		port.Value = *p.From
	}

	return port
}

type PortList struct {
	Labels []string
	Ports  []Port
}

func NewPortList(ports []Port) *PortList {
	pl := &PortList{}
	for _, p := range ports {
		pl.Labels = append(pl.Labels, p.Label)
		pl.Ports = append(pl.Ports, p)
	}
	return pl
}

func BuildDynamicPortList(portTargets []int) *PortList {
	ports := []Port{}
	for ind, p := range portTargets {
		ports = append(ports, *NewPort(fmt.Sprintf("http-%d", ind), p, nil))
	}
	return NewPortList(ports)
}

func (pl *PortList) GetDynamic() []nomadStructs.Port {
	dynamic := []nomadStructs.Port{}
	for _, p := range pl.Ports {
		if p.From == nil {
			dynamic = append(dynamic, p.ToNomadPort())
		}
	}
	return dynamic
}

func (pl *PortList) GetStatic() []nomadStructs.Port {
	static := []nomadStructs.Port{}
	for _, p := range pl.Ports {
		if p.From != nil {
			static = append(static, p.ToNomadPort())
		}
	}
	return static
}
