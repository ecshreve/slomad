package registry

import (
	"bytes"
	"encoding/gob"
	"fmt"

	nomadApi "github.com/hashicorp/nomad/api"
	nomadStructs "github.com/hashicorp/nomad/nomad/structs"
	"github.com/kr/pretty"
	"github.com/samsarahq/go/oops"
	log "github.com/sirupsen/logrus"
)

// convertJob converts a nomad struct job to a nomad api job.
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

// planApiJob creates a nomad api client, and runs a plan
// for the given job, printing the output diff.
func planApiJob(job *nomadApi.Job) error {
	nomadConfig := nomadApi.DefaultConfig()
	nomadClient, err := nomadApi.NewClient(nomadConfig)
	if err != nil {
		return oops.Wrapf(err, "unable to create nomad api client")
	}

	planResp, _, nomadErr := nomadClient.Jobs().Plan(job, true, nil)
	if nomadErr != nil {
		log.Errorf("Error submitting job: %s", nomadErr)
		return fmt.Errorf(fmt.Sprintf("Error submitting job: %s", nomadErr))
	}

	log.Infof("Sucessfully planned nomad job. %v\n", job.Name)
	pretty.Print(planResp.Diff)
	return nil
}

// submitApiJob creates a nomad api client, and submits the job to nomad.
func submitApiJob(job *nomadApi.Job) error {
	nomadConfig := nomadApi.DefaultConfig()
	nomadClient, err := nomadApi.NewClient(nomadConfig)
	if err != nil {
		return oops.Wrapf(err, "unable to create nomad api client")
	}

	evalId, _, nomadErr := nomadClient.Jobs().Register(job, nil)
	if nomadErr != nil {
		return oops.Wrapf(err, "error submitting job: %+v", job)
	}

	log.Infof("Sucessfully submitted nomad job. Eval id: %v\n", evalId)
	return nil
}

// // Plan creates a new API job and runs a plan on it.
// func (j *Job) Plan(force bool) error {
// 	apiJob, err := j.CreateApiJob(force)
// 	if err != nil {
// 		return oops.Wrapf(err, "error crating api job for job: %+v", j)
// 	}

// 	if err = planApiJob(apiJob); err != nil {
// 		return oops.Wrapf(err, "error planning api job")
// 	}

// 	return nil
// }

// // Deploy creates a new ApiJob and submits it to the nomad API.
// func (j *Job) Deploy(force bool) error {
// 	apiJob, err := j.CreateApiJob(force)
// 	if err != nil {
// 		return oops.Wrapf(err, "error creating api job for job: %+v", j)
// 	}

// 	if err = submitApiJob(apiJob); err != nil {
// 		return oops.Wrapf(err, "error submitting api job")
// 	}

// 	return nil
// }
