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

type ResourceValue int

const (
	TINY    ResourceValue = 128
	SMALL   ResourceValue = 256
	DEFAULT ResourceValue = 512
	LARGE   ResourceValue = 1024
	XLARGE  ResourceValue = 2048
	XXLARGE ResourceValue = 4096
)

// TaskResource is a struct that represents the CPU and MEM resources for a task.
type TaskResource struct {
	CPU ResourceValue
	MEM ResourceValue
}

var (
	DEFAULT_TASK = TaskResource{CPU: DEFAULT, MEM: DEFAULT}
	TINY_TASK    = TaskResource{CPU: TINY, MEM: TINY}
	SMALL_TASK   = TaskResource{CPU: SMALL, MEM: SMALL}
	LARGE_TASK   = TaskResource{CPU: LARGE, MEM: LARGE}
	XLARGE_TASK  = TaskResource{CPU: XLARGE, MEM: XLARGE}
	MEM_TASK     = TaskResource{CPU: SMALL, MEM: LARGE}
	COMPUTE_TASK = TaskResource{CPU: LARGE, MEM: SMALL}
	PLEX_TASK    = TaskResource{CPU: LARGE, MEM: XXLARGE}
)
