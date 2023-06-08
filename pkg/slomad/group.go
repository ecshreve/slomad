package slomad

import (
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

// getReschedulePolicy returns a nomad reschedule policy for a given job.
func getReschedulePolicy(jt JobType) *nomadStructs.ReschedulePolicy {
	if jt == SERVICE || jt == STORAGE_CONTROLLER {
		return &nomadStructs.DefaultServiceJobReschedulePolicy
	}
	return nil
}

// getDisk returns a nomad disk struct with a default size for a given job.
func getDisk() *nomadStructs.EphemeralDisk {
	return &nomadStructs.EphemeralDisk{
		SizeMB: 500,
	}
}
