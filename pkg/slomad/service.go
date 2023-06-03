package slomad

import (
	"fmt"
	"time"

	nomadStructs "github.com/hashicorp/nomad/nomad/structs"
)

func GetService(taskName string, portLabel string) *nomadStructs.Service {
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

func GetTask(j *JJJob) *nomadStructs.Task {
	portLabels := ExtractLabels(j.Ports)
	cfg := map[string]interface{}{
		"image": j.Image,
		"args":  j.Args,
		"ports": portLabels,
	}

	services := []*nomadStructs.Service{}
	for _, pl := range portLabels {
		services = append(services, GetService(j.Name, pl))
	}

	return &nomadStructs.Task{
		Name:      j.Name,
		Driver:    "docker",
		Config:    cfg,
		Resources: j.Shape.ToNomadResource(),
		Services:  services,
		LogConfig: nomadStructs.DefaultLogConfig(),
	}
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

func getDisk() *nomadStructs.EphemeralDisk {
	return &nomadStructs.EphemeralDisk{
		SizeMB: 500,
	}
}

func GetGroup(j *JJJob) *nomadStructs.TaskGroup {
	return &nomadStructs.TaskGroup{
		Name:             j.Name,
		Count:            1,
		Tasks:            []*nomadStructs.Task{GetTask(j)},
		RestartPolicy:    nomadStructs.NewRestartPolicy("service"),
		ReschedulePolicy: &nomadStructs.DefaultServiceJobReschedulePolicy,
		EphemeralDisk:    getDisk(),
		Networks:         getNetworks(j.Ports),
	}
}
