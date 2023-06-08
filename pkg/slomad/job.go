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
	Ports     []*Port
	Volumes   []Volume
	Args      []string
	Env       map[string]string
	User      string
	Templates map[string]string
}

// GetNomadApiJob returns a nomadApi Job for the given slomad.Job.
func (j *Job) GetNomadApiJob(force bool) (*nomadApi.Job, error) {
	job := &nomadStructs.Job{
		Priority:    50,
		Namespace:   "default",
		Region:      "global",
		Datacenters: []string{"dcs"},

		ID:          j.Name,
		Name:        j.Name,
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

	if err := job.Validate(); err != nil {
		return nil, oops.Wrapf(err, "Nomad job validation failed")
	}

	return convertJob(job)
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

// getConstraint returns a nomad constraint for a given deploy target.
func getConstraint(dt DeployTarget) *nomadStructs.Constraint {
	return &nomadStructs.Constraint{
		LTarget: "${attr.unique.hostname}",
		RTarget: DeployTargetRegex[dt],
		Operand: "regexp",
	}
}
