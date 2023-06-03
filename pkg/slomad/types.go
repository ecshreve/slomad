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
)

// DeployTargetRegex is a map of DeployTarget to a regex string that matches
// the target's hostname.
var DeployTargetRegex = map[DeployTarget]string{
	ALL:      "^.*$",
	SERVER:   "^server-[0-9]+$",
	WORKER:   "^worker-[0-9]+$",
	CODERBOX: "^coderbox$",
	DEVBOX:   "^devbox$",
}

// JobType is an enum that represents the type of a job.
type JobType int

const (
	UNKNOWN_JOB_TYPE JobType = iota
	SERVICE
	SYSTEM
	BATCH
)

// String implements the Stringer interface for JobType.
func (jt JobType) String() string {
	return [...]string{"UNKNOWN", "service", "system", "batch"}[jt]
}
