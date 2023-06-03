package slomad

import (
	"bytes"
	"encoding/gob"
	"fmt"

	"github.com/google/uuid"
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
type App struct {
	Name       string
	Type       JobType
	Shape      TaskResource
	Constraint string
	Image      string
	Args       []string
	Ports      []*Port
	Env        map[string]string
	User       string
	Storage    string
	Volumes    []Volume
	Templates  map[string]string
}

// GetConstraint returns a nomad constraint for a given job.
func getConstraint(j *App) *nomadStructs.Constraint {
	return &nomadStructs.Constraint{
		LTarget: "${attr.unique.hostname}",
		RTarget: j.Constraint,
		Operand: "regexp",
	}
}

// ToNomadJob converts a JJJob to a Nomad Job
func (j *App) ToNomadJob(force bool) (*nomadStructs.Job, *nomadApi.Job, error) {
	job := &nomadStructs.Job{
		Priority:    50,
		Namespace:   "default",
		Region:      "global",
		Datacenters: []string{"dcs"},

		ID:   j.Name,
		Name: j.Name,

		Type:        j.Type.String(),
		TaskGroups:  []*nomadStructs.TaskGroup{GetGroup(j)},
		Constraints: []*nomadStructs.Constraint{getConstraint(j)},
	}

	// Writing a new uuid to this field ensures Nomad will create a new
	// version of the job.
	if force {
		job.Meta = map[string]string{
			"run_uuid": uuid.NewString(),
		}
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

// convertJob converts a Nomad Job to a Nomad API Job.
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
	StorageParams
}

type StorageParams struct {
	Storage *string
	Volumes []Volume
}

type TaskConfigParams struct {
	Args      []string
	Ports     []*Port
	Shape     TaskResource
	Env       map[string]string
	User      string
	Templates map[string]string
}

func NewServiceJob(params JobParams) *App {
	return &App{
		Name:       params.Name,
		Image:      fmt.Sprintf("reg.slab.lan:5000/%s", params.Name),
		Args:       params.Args,
		Ports:      params.Ports,
		Type:       SERVICE,
		Shape:      params.Shape,
		Constraint: DeployTargetRegex[params.Target],
		Env:        params.Env,
		User:       params.User,
		Volumes:    params.Volumes,
		Storage:    StringValOr(params.Storage, ""),
		Templates:  params.Templates,
	}
}
