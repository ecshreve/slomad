package slomad

import (
	"fmt"
	"time"

	nomadStructs "github.com/hashicorp/nomad/nomad/structs"
)

// getService returns a nomad service struct for a given task.
func getService(taskName string, portLabel string) *nomadStructs.Service {
	if taskName == "storage-controller" || taskName == "storage-node" {
		return nil
	}

	return &nomadStructs.Service{
		Name:      taskName,
		PortLabel: portLabel,
		TaskName:  taskName,
		Tags:      []string{"traefik.enable=true"},
		Checks: []*nomadStructs.ServiceCheck{
			{
				Name:          fmt.Sprintf("%s -- %s = tcp check", taskName, portLabel),
				Type:          nomadStructs.ServiceCheckTCP,
				Interval:      10 * time.Second,
				Timeout:       2 * time.Second,
				InitialStatus: "passing",
			},
		},
		Provider: "consul",
	}
}

// getServices returns a list of services for a given job.
func getServices(taskName string, portLabels []string) []*nomadStructs.Service {
	services := []*nomadStructs.Service{}
	for _, pl := range portLabels {
		srvc := getService(taskName, pl)
		if srvc == nil {
			continue
		}
		services = append(services, srvc)
	}
	return services
}

func getServiceWithTags(taskName, portLabel string, tags []string) *nomadStructs.Service {
	srvc := getService(taskName, portLabel)

	srvc.Tags = tags
	return srvc
}

func getGroupService(name, portLabel string) *nomadStructs.Service {
	srvc := getService(name, portLabel)
	srvc.Tags = nil
	srvc.TaskName = ""
	return srvc
}
