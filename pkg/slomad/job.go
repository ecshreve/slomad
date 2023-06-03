package slomad

import (
	"bytes"
	"encoding/gob"
	"fmt"

	nomadApi "github.com/hashicorp/nomad/api"
	nomadStructs "github.com/hashicorp/nomad/nomad/structs"
	log "github.com/sirupsen/logrus"
)

type JobType int

const (
	UNKNOWN_JobType JobType = iota
	SERVICE
	SYSTEM
)

func (jt JobType) String() string {
	return [...]string{"UNKNOWN", "service", "system"}[jt]
}

type CommonArgs struct {
	Driver     string
	Constraint string
	Count      int
	Priority   int
}

type Job struct {
	Name      string
	Image     *string
	JobType   string
	Storage   *string
	User      *string
	Size      map[string]int
	Ports     []Port
	Caps      []string
	Args      []string
	Env       map[string]string
	Templates map[string]string
	Volumes   map[string]string
	Mounts    map[string]string
	CommonArgs
}

// ++++++++++++

// JJJob is a new type of job
type JJJob struct {
	Name       string
	Type       JobType
	Shape      TaskResource
	Constraint string
	Image      string
	Args       []string
	Ports      []*Port
}

// ToNomadJob converts a JJJob to a Nomad Job
func (j *JJJob) ToNomadJob() (*nomadStructs.Job, *nomadApi.Job, error) {
	job := &nomadStructs.Job{
		ID:          j.Name,
		Name:        j.Name,
		Region:      "global",
		Priority:    50,
		Datacenters: []string{"dcs"},
		Type:        j.Type.String(),
		TaskGroups:  []*nomadStructs.TaskGroup{GetGroup(j)},
		Namespace:   "default",
		Constraints: []*nomadStructs.Constraint{
			{
				LTarget: "${attr.unique.hostname}",
				RTarget: j.Constraint,
				Operand: "regexp",
			},
		},
	}

	if err := job.Validate(); err != nil {
		log.Errorf("Nomad job validation failed. Error: %s\n", err)
		return job, nil, err
	}

	apiJob, err := convertJob(job)
	if err != nil {
		log.Errorf("Failed to convert nomad job in api call. Error: %s\n", err)
		return job, apiJob, err
	}

	return job, apiJob, nil
}

func convertJob(in *nomadStructs.Job) (*nomadApi.Job, error) {
	gob.Register([]map[string]interface{}{})
	gob.Register([]interface{}{})

	var apiJob *nomadApi.Job
	buf := new(bytes.Buffer)
	if err := gob.NewEncoder(buf).Encode(in); err != nil {
		return nil, err
	}
	if err := gob.NewDecoder(buf).Decode(&apiJob); err != nil {
		return nil, err
	}

	return apiJob, nil
}

type JobParams struct {
	Name   string
	Target DeployTarget
	TaskConfigParams
}

type TaskConfigParams struct {
	Args  []string
	Ports []PortParams
	Shape TaskResource
}

func NewServiceJob(params JobParams) *JJJob {
	return &JJJob{
		Name:       params.Name,
		Image:      fmt.Sprintf("reg.slab.lan:5000/%s", params.Name),
		Args:       params.Args,
		Ports:      NewPorts(params.Ports),
		Type:       SERVICE,
		Shape:      params.Shape,
		Constraint: DeployTargetRegex[params.Target],
	}
}
