package slomad

import (
	"fmt"
	"time"

	nomadStructs "github.com/hashicorp/nomad/nomad/structs"
)

func GetGroup(j *App) *nomadStructs.TaskGroup {
	return &nomadStructs.TaskGroup{
		Name:             j.Name,
		Count:            1,
		Tasks:            []*nomadStructs.Task{GetTask(j)},
		RestartPolicy:    nomadStructs.NewRestartPolicy("service"),
		ReschedulePolicy: &nomadStructs.DefaultServiceJobReschedulePolicy,
		EphemeralDisk:    getDisk(),
		Networks:         getNetworks(j.Ports),
		Volumes:          getNomadVolumes(j.Storage),
	}
}

func GetTask(j *App) *nomadStructs.Task {
	return &nomadStructs.Task{
		Name:         j.Name,
		Driver:       "docker",
		Config:       getConfig(j),
		Resources:    j.Shape.ToNomadResource(),
		Services:     getServices(j.Name, ExtractLabels(j.Ports)),
		LogConfig:    nomadStructs.DefaultLogConfig(),
		Env:          j.Env,
		User:         j.User,
		VolumeMounts: getMounts(j.Volumes),
	}
}

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

// getServices returns a list of services for a given job.
func getServices(taskName string, portLabels []string) []*nomadStructs.Service {
	services := []*nomadStructs.Service{}
	for _, pl := range portLabels {
		services = append(services, GetService(taskName, pl))
	}
	return services
}

// getConfig returns a nomad config struct for a given job.
func getConfig(j *App) map[string]interface{} {
	portLabels := ExtractLabels(j.Ports)
	config := map[string]interface{}{
		"image": j.Image,
		"args":  j.Args,
		"ports": portLabels,
	}

	vols := toVolumeStrings(j.Volumes)
	if len(vols) > 0 {
		config["volumes"] = vols
	}

	return config
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
