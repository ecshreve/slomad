package slomad

import (
	nomadStructs "github.com/hashicorp/nomad/nomad/structs"
	log "github.com/sirupsen/logrus"
)

// getGroup returns a nomad task group struct for a given job.
func getGroup(j *Job) *nomadStructs.TaskGroup {
	networkMode := ""
	disk := 500
	if j.Name == "traefik" {
		networkMode = "host"
		disk = 256
	}

	tg := &nomadStructs.TaskGroup{
		Name:             j.Name,
		Count:            1,
		Tasks:            []*nomadStructs.Task{getTask(j)},
		RestartPolicy:    nomadStructs.NewRestartPolicy(j.Type.String()),
		ReschedulePolicy: getReschedulePolicy(j.Type),
		EphemeralDisk:    getDisk(disk),
		Networks:         getNetworks(j.Ports, networkMode),
		Volumes:          getNomadVolumeReq(j.Volumes),
	}

	if j.Name == "traefik" {
		tg.Services = []*nomadStructs.Service{
			getGroupService("traefik-web", "web"),
		}
	}

	return tg
}

// getReschedulePolicy returns a nomad reschedule policy for a given job.
func getReschedulePolicy(jt JobType) *nomadStructs.ReschedulePolicy {
	if jt == SERVICE || jt == STORAGE_CONTROLLER {
		return &nomadStructs.DefaultServiceJobReschedulePolicy
	}
	return nil
}

// getDisk returns a nomad disk struct with a default size for a given job.
func getDisk(mb int) *nomadStructs.EphemeralDisk {
	if mb <= 0 || mb > 2000 {
		log.Infof("Invalid disk size %d, using default of 500MB", mb)
		mb = 500
	}
	return &nomadStructs.EphemeralDisk{
		SizeMB: mb,
	}
}
