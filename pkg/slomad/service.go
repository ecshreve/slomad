package slomad

import (
	"fmt"
	"time"

	nomadStructs "github.com/hashicorp/nomad/nomad/structs"
)

func GetGroup(j *Job) *nomadStructs.TaskGroup {
	return &nomadStructs.TaskGroup{
		Name:             j.Name,
		Count:            1,
		Tasks:            []*nomadStructs.Task{GetTask(j)},
		RestartPolicy:    nomadStructs.NewRestartPolicy(j.Type.String()),
		ReschedulePolicy: getReschedulePolicy(j),
		EphemeralDisk:    getDisk(),
		Networks:         getNetworks(j.Ports),
		Volumes:          getNomadVolumes(j.Storage),
	}
}

func GetTask(j *Job) *nomadStructs.Task {
	return &nomadStructs.Task{
		Name:            j.Name,
		Driver:          "docker",
		Config:          getConfig(j),
		Resources:       getResource(j.Shape),
		Services:        getServices(j.Name, extractLabels(j.Ports)),
		LogConfig:       nomadStructs.DefaultLogConfig(),
		Env:             j.Env,
		User:            j.User,
		VolumeMounts:    getMounts(j.Volumes),
		Templates:       getTemplates(j.Templates),
		CSIPluginConfig: getCSIPluginConfig(j),
	}
}

func GetService(taskName string, portLabel string) *nomadStructs.Service {
	if taskName == "storage-controller" {
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

func getTemplates(templates map[string]string) []*nomadStructs.Template {
	if templates == nil {
		return nil
	}

	nt := []*nomadStructs.Template{}
	for tmplname, tmpl := range templates {
		nt = append(nt, &nomadStructs.Template{
			EmbeddedTmpl: tmpl,
			DestPath:     fmt.Sprintf("local/config/%s", tmplname),
			ChangeMode:   "signal",
			ChangeSignal: "SIGHUP",
		})
	}

	return nt
}

// getServices returns a list of services for a given job.
func getServices(taskName string, portLabels []string) []*nomadStructs.Service {
	services := []*nomadStructs.Service{}
	for _, pl := range portLabels {
		srvc := GetService(taskName, pl)
		if srvc == nil {
			continue
		}
		services = append(services, srvc)
	}
	return services
}

// getConfig returns a nomad config struct for a given job.
func getConfig(j *Job) map[string]interface{} {
	portLabels := extractLabels(j.Ports)
	config := map[string]interface{}{
		"image": j.Image,
		"args":  j.Args,
		"ports": portLabels,
	}

	vols := getVolumeStrings(j.Volumes)
	if len(vols) > 0 {
		config["volumes"] = vols
	}

	if j.Storage == "controller" || j.Storage == "node" {
		config["privileged"] = true
		config["network_mode"] = "host"
	}

	return config
}

func getDisk() *nomadStructs.EphemeralDisk {
	return &nomadStructs.EphemeralDisk{
		SizeMB: 500,
	}
}

// GetConstraint returns a nomad constraint for a given job.
func getConstraint(j *Job) *nomadStructs.Constraint {
	return &nomadStructs.Constraint{
		LTarget: "${attr.unique.hostname}",
		RTarget: j.Constraint,
		Operand: "regexp",
	}
}

func getReschedulePolicy(j *Job) *nomadStructs.ReschedulePolicy {
	if j.Type == SERVICE {
		return &nomadStructs.DefaultServiceJobReschedulePolicy
	}
	return nil
}
