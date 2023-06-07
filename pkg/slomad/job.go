package slomad

import (
	"bytes"
	"encoding/gob"

	"github.com/google/uuid"
	nomadApi "github.com/hashicorp/nomad/api"
	nomadStructs "github.com/hashicorp/nomad/nomad/structs"
	"github.com/samsarahq/go/oops"
)

// Job is a struct that represents a Nomad Job.
type Job struct {
	Name      string
	Type      JobType
	Shape     TaskResource
	Target    DeployTarget
	Args      []string
	Ports     []*Port
	Env       map[string]string
	User      string
	Volumes   []Volume
	Templates map[string]string
}

// ToNomadJob converts a Job to a it's Nomad representations.
func (j *Job) ToNomadJob(force bool) (*nomadStructs.Job, *nomadApi.Job, error) {
	job := &nomadStructs.Job{
		Priority:    50,
		Namespace:   "default",
		Region:      "global",
		Datacenters: []string{"dcs"},

		ID:   j.Name,
		Name: j.Name,

		Type:        j.Type.String(),
		TaskGroups:  []*nomadStructs.TaskGroup{getGroup(j)},
		Constraints: []*nomadStructs.Constraint{getConstraint(j.Target)},
	}

	// Writing a new uuid to this field ensures Nomad will create a new
	// version of the job.
	if force {
		job.Meta = map[string]string{
			"run_uuid": uuid.NewString(),
		}
	}

	apiJob, err := convertJob(job)
	if err != nil {
		return nil, nil, oops.Wrapf(err, "Failed to convert nomad job for api call")
	}

	return job, apiJob, nil
}

// convertJob converts a Nomad Job to a Nomad API Job.
func convertJob(in *nomadStructs.Job) (*nomadApi.Job, error) {
	if err := in.Validate(); err != nil {
		return nil, oops.Wrapf(err, "Nomad job validation failed")
	}

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

// JobParams is a struct that represents the parameters for creating a Job.
type JobParams struct {
	Name   string
	Type   JobType
	Target DeployTarget
	TaskConfigParams
	StorageParams
}

// StorageParams is a struct that represents the parameters for creating a
// storage associated with a Job.
type StorageParams struct {
	Volumes []Volume
}

// TaskConfigParams is a struct that represents the parameters for creating a
// TaskConfig associated with a Job.
type TaskConfigParams struct {
	Args      []string
	Ports     []*Port
	Shape     TaskResource
	Env       map[string]string
	User      string
	Templates map[string]string
}

// NewAppJob creates a new Job for an application.
func NewJob(params JobParams) *Job {
	return &Job{
		Name:      params.Name,
		Args:      params.Args,
		Ports:     params.Ports,
		Type:      params.Type,
		Target:    params.Target,
		Shape:     params.Shape,
		Env:       params.Env,
		User:      params.User,
		Volumes:   params.Volumes,
		Templates: params.Templates,
	}
}
