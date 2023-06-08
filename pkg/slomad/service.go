package slomad

import (
	"fmt"
	"time"

	nomadStructs "github.com/hashicorp/nomad/nomad/structs"
)

// getGroup returns a nomad task group struct for a given job.
func getGroup(j *Job) *nomadStructs.TaskGroup {
	return &nomadStructs.TaskGroup{
		Name:             j.Name,
		Count:            1,
		Tasks:            []*nomadStructs.Task{getTask(j)},
		RestartPolicy:    nomadStructs.NewRestartPolicy(j.Type.String()),
		ReschedulePolicy: getReschedulePolicy(j.Type),
		EphemeralDisk:    getDisk(),
		Networks:         getNetworks(j.Ports),
		Volumes:          getNomadVolumeReq(j.Volumes),
	}
}

// getTask returns a nomad task struct for a given job.
func getTask(j *Job) *nomadStructs.Task {
	return &nomadStructs.Task{
		Name:            j.Name,
		Driver:          "docker",
		Config:          getConfig(j.Name, j.Type, j.Args, j.Ports, j.Volumes),
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

func getTemplates(templates map[string]string) []*nomadStructs.Template {
	if templates == nil || len(templates) == 0 {
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
		srvc := getService(taskName, pl)
		if srvc == nil {
			continue
		}
		services = append(services, srvc)
	}
	return services
}

// getConfig returns a nomad config struct for a given job.
func getConfig(name string, jt JobType, args []string, ports []*Port, vols []Volume) map[string]interface{} {
	config := map[string]interface{}{
		"image": fmt.Sprintf("reg.slab.lan:5000/%s", name),
		"args":  args,
		"ports": extractLabels(ports),
	}

	volStrings := getVolumeStrings(vols)
	if len(volStrings) > 0 {
		config["volumes"] = volStrings
	}

	if jt == STORAGE_CONTROLLER || jt == STORAGE_NODE {
		config["privileged"] = true
		config["network_mode"] = "host"
		config["image"] = "reg.slab.lan:5000/csi-nfs-plugin"
	}

	if name == "plex" {
		config["network_mode"] = "host"
	}

	return config
}

// getDisk returns a nomad disk struct with a default size for a given job.
func getDisk() *nomadStructs.EphemeralDisk {
	return &nomadStructs.EphemeralDisk{
		SizeMB: 500,
	}
}

// getConstraint returns a nomad constraint for a given deploy target.
func getConstraint(dt DeployTarget) *nomadStructs.Constraint {
	return &nomadStructs.Constraint{
		LTarget: "${attr.unique.hostname}",
		RTarget: DeployTargetRegex[dt],
		Operand: "regexp",
	}
}

// getReschedulePolicy returns a nomad reschedule policy for a given job.
func getReschedulePolicy(jt JobType) *nomadStructs.ReschedulePolicy {
	if jt == SERVICE || jt == STORAGE_CONTROLLER {
		return &nomadStructs.DefaultServiceJobReschedulePolicy
	}
	return nil
}
