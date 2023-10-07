package slomad

import (
	"fmt"

	nomadStructs "github.com/hashicorp/nomad/nomad/structs"
)

// getTask returns a nomad task struct for a given job.
func getTask(j *Job) *nomadStructs.Task {
	labels := extractLabels(j.Ports)
	srvcs := getServices(j.Name, labels)

	tags := []string{}
	if j.TaskServiceTags != nil {
		tags = j.TaskServiceTags[j.Name]
		for _, srvc := range srvcs {
			srvc.Tags = tags
		}
	}

	return &nomadStructs.Task{
		Name:            j.Name,
		Driver:          "docker",
		Config:          getConfig(j.Name, j.Type, j.Args, j.Ports, j.Volumes),
		Resources:       getResource(j.Shape),
		Services:        srvcs,
		LogConfig:       nomadStructs.DefaultLogConfig(),
		Env:             j.Env,
		User:            j.User,
		VolumeMounts:    getMounts(j.Volumes),
		Templates:       getTemplates(j.Templates),
		CSIPluginConfig: getCSIPluginConfig(j),
	}
}

// getResource is a helper function to convert a TaskResource to a Nomad Resources struct.
func getResource(tr TaskResource) *nomadStructs.Resources {
	return &nomadStructs.Resources{
		CPU:      int(tr.CPU),
		MemoryMB: int(tr.MEM),
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

// getConfig returns a nomad config struct for a given job.
func getConfig(name string, jt JobType, args []string, ports []*Port, vols []Volume) map[string]interface{} {
	config := map[string]interface{}{
		"image": fmt.Sprintf("registry.slab.lan:5000/%s", name),
		"args":  args,
		"ports": extractLabels(ports),
	}

	if jt == STORAGE_CONTROLLER {
		config["image"] = "registry.gitlab.com/rocketduck/csi-plugin-nfs:0.7.0"
		config["privileged"] = true
		config["network_mode"] = "host"
	}

	volStrings := getVolumeStrings(vols)
	if len(volStrings) > 0 {
		config["volumes"] = volStrings
	}

	return config
}
