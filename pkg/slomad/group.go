package slomad

import (
	nomadStructs "github.com/hashicorp/nomad/nomad/structs"
	log "github.com/sirupsen/logrus"
)

// getGroup returns a nomad task group struct for a given job.
func getGroup(j *Job) *nomadStructs.TaskGroup {
	tg := &nomadStructs.TaskGroup{
		Name:             j.Name,
		Count:            1,
		Tasks:            []*nomadStructs.Task{getTask(j)},
		RestartPolicy:    nomadStructs.NewRestartPolicy(j.Type.String()),
		ReschedulePolicy: getReschedulePolicy(j.Type),
		EphemeralDisk:    getDisk(0), // TODO: make this configurable
		Networks:         getNetworks(j.Ports, j.Priv),
		Volumes:          getNomadVolumeReq(j.Volumes),
	}

	if j.GroupServices != nil {
		tg.Services = getGroupServices(j.GroupServices)
	}

	return tg
}

// getReschedulePolicy returns a nomad reschedule policy for a given job.
func getReschedulePolicy(jt JobType) *nomadStructs.ReschedulePolicy {
	if jt == SERVICE || jt == STORAGE_CONTROLLER {
		return &nomadStructs.DefaultServiceJobReschedulePolicy
	}

	if jt == BATCH {
		return &nomadStructs.DefaultBatchJobReschedulePolicy
	}

	return nil
}

// getDisk returns a nomad disk struct with a default size for a given job.
func getDisk(mb int) *nomadStructs.EphemeralDisk {
	if mb <= 0 || mb > 2000 {
		log.Debugf("Invalid disk size %d, using default of 500MB", mb)
		mb = 500
	}
	return &nomadStructs.EphemeralDisk{
		SizeMB: mb,
	}
}
