package slomad

type DeployTarget int

const (
	UNKNOWN DeployTarget = iota
	ALL
	SERVER
	WORKER
	CODERBOX
	DEVBOX
)

var DeployTargetRegex = map[DeployTarget]string{
	ALL:      "^.*$",
	SERVER:   "^server-[0-9]+$",
	WORKER:   "^worker-[0-9]+$",
	CODERBOX: "^coderbox$",
	DEVBOX:   "^devbox$",
}

type JobType int

const (
	UNKNOWN_JobType JobType = iota
	SERVICE
	SYSTEM
	BATCH
)

func (jt JobType) String() string {
	return [...]string{"UNKNOWN", "service", "system", "batch"}[jt]
}
