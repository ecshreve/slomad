package slomad

import "fmt"

type DeployTarget int

const (
	UNKNOWN DeployTarget = iota
	SERVER
	WORKER
	CODERBOX
	DEVBOX
)

type Targ interface {
	GetRegex() string
}

func (d DeployTarget) GetRegex() string {
	return "^worker-[0-9]+$"
}

func NewJob(name string, d DeployTarget) *Job {
	// dockerImage := getDockerImageString(name)
	return nil
}

func getDockerImageString(name string) *string {
	imageStr := fmt.Sprintf("registry.slab.lan:5000/%s:custom", name)
	return &imageStr
}
