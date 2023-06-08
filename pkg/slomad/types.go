package slomad

// DeployTarget is an enum that represents the target of a job.
type DeployTarget int

const (
	UNKNOWN_DEPLOY_TARGET DeployTarget = iota
	ALL
	SERVER
	WORKER
	CODERBOX
	DEVBOX
	PLEXBOX
)

// DeployTargetRegex is a map of DeployTarget to a regex string that matches
// the target's hostname.
var DeployTargetRegex = map[DeployTarget]string{
	ALL:      "^.*$",
	SERVER:   "^server-[0-9]+$",
	WORKER:   "^worker-[0-9]+$",
	CODERBOX: "^coderbox$",
	DEVBOX:   "^devbox$",
	PLEXBOX:  "^plexbox$",
}

// JobType is an enum that represents the type of a job.
type JobType int

const (
	UNKNOWN_JOB_TYPE JobType = iota
	SERVICE
	SYSTEM
	BATCH
	STORAGE_CONTROLLER
	STORAGE_NODE
)

// String implements the Stringer interface for JobType.
func (jt JobType) String() string {
	return [...]string{"UNKNOWN", "service", "system", "batch", "service", "system"}[jt]
}
