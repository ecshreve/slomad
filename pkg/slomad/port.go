package slomad

import (
	nomadStructs "github.com/hashicorp/nomad/nomad/structs"
)

// Port is a struct that represents a network port for a task.
type Port struct {
	Label  string
	To     int
	From   int
	Static bool
}

func toNomadPort(p *Port) nomadStructs.Port {
	return nomadStructs.Port{
		Label: p.Label,
		To:    p.To,
		Value: p.From,
	}
}

// basicPort returns a basic port struct with a default label.
func basicPort(val int) *Port {
	return &Port{
		Label:  "http",
		To:     val,
		From:   0,
		Static: false,
	}
}

// BasicPortConfig returns a list with a single Port element with a default label.
func BasicPortConfig(val int) []*Port {
	return []*Port{basicPort(val)}
}

// toNomadPortMap converts a list of Ports to a map of static and dynamic ports.
func toNomadPortMap(ports []*Port) map[string][]nomadStructs.Port {
	stat := []nomadStructs.Port{}
	dynm := []nomadStructs.Port{}

	for _, p := range ports {
		np := toNomadPort(p)
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

// extractLabels returns a list of labels from a list of Ports.
func extractLabels(ports []*Port) []string {
	labels := []string{}
	for _, p := range ports {
		labels = append(labels, p.Label)
	}
	return labels
}

// getNetworks converts a list of Ports to a list of Nomad NetworkResources.
func getNetworks(ports []*Port, mode string) []*nomadStructs.NetworkResource {
	portMap := toNomadPortMap(ports)
	return []*nomadStructs.NetworkResource{
		{
			Mode:          mode,
			ReservedPorts: portMap["static"],
			DynamicPorts:  portMap["dynamic"],
		},
	}
}
