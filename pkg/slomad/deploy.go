package slomad

import (
	"fmt"
	"os"

	nomadApi "github.com/hashicorp/nomad/api"
	"github.com/samsarahq/go/oops"
	log "github.com/sirupsen/logrus"
)

// planApiJob creates a nomad api client, and runs a plan
// for the given job, printing the output diff.
func planApiJob(job *nomadApi.Job) error {
	nomadConfig := nomadApi.DefaultConfig()
	nomadConfig.Address = os.Getenv("NOMAD_TARGET")
	nomadClient, err := nomadApi.NewClient(nomadConfig)
	if err != nil {
		return oops.Wrapf(err, "unable to create nomad api client")
	}

	planResp, _, nomadErr := nomadClient.Jobs().Plan(job, true, nil)
	if nomadErr != nil {
		log.Errorf("Error planning job: %s", nomadErr)
		return fmt.Errorf(fmt.Sprintf("Error planning job: %s", nomadErr))
	}

	log.Infof("Sucessfully planned nomad job %s - %+v\n", *job.Name, planResp.Annotations.DesiredTGUpdates[*job.Name])
	return nil
}

// submitApiJob creates a nomad api client, and submits the job to nomad.
//
// TODO: move client creation to a helper function
func submitApiJob(job *nomadApi.Job) error {
	nomadConfig := nomadApi.DefaultConfig()
	nomadConfig.Address = os.Getenv("NOMAD_TARGET")
	nomadClient, err := nomadApi.NewClient(nomadConfig)
	if err != nil {
		return oops.Wrapf(err, "unable to create nomad api client")
	}

	_, _, nomadErr := nomadClient.Jobs().Register(job, nil)
	if nomadErr != nil {
		return oops.Wrapf(err, "error submitting job: %+v", job)
	}

	log.Infof("Sucessfully submitted nomad job %s\n", *job.Name)
	return nil
}

// Plan creates a new API job and runs a plan on it.
func (j *Job) Plan(force bool) error {
	_, aj, err := j.ToNomadJob(force)
	if err != nil {
		return oops.Wrapf(err, "error creating api job for job: %+v", j)
	}

	if err = planApiJob(aj); err != nil {
		return oops.Wrapf(err, "error planning api job")
	}

	return nil
}

func (j *Job) Deploy(force bool) error {
	_, aj, err := j.ToNomadJob(force)
	if err != nil {
		return oops.Wrapf(err, "error creating api job for job: %+v", j)
	}

	if err = submitApiJob(aj); err != nil {
		return oops.Wrapf(err, "error submitting api job")
	}

	return nil
}
