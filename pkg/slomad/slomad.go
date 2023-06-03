package slomad

type DeployTarget int

const (
	UNKNOWN DeployTarget = iota
	SERVER
	WORKER
	CODERBOX
	DEVBOX
)

var DeployTargetRegex = map[DeployTarget]string{
	SERVER:   "^server-[0-9]+$",
	WORKER:   "^worker-[0-9]+$",
	CODERBOX: "^coderbox$",
	DEVBOX:   "^devbox$",
}

type TaskParams struct {
	Driver     string
	Constraint string
	Count      int
	Priority   int
}

func (dt DeployTarget) DefaultTaskParams() TaskParams {
	return TaskParams{
		Driver:     "docker",
		Constraint: DeployTargetRegex[dt],
		Count:      1,
		Priority:   50,
	}
}
