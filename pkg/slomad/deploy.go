package slomad

import (
	"fmt"
	"os"

	nomadApi "github.com/hashicorp/nomad/api"
	"github.com/samsarahq/go/oops"
	log "github.com/sirupsen/logrus"
)

func newNomadClient() (*nomadApi.Client, error) {
	nomadConfig := nomadApi.DefaultConfig()
	nomadConfig.Address = os.Getenv("NOMAD_TARGET")
	nomadClient, err := nomadApi.NewClient(nomadConfig)
	if err != nil {
		return nil, oops.Wrapf(err, "unable to create nomad api client")
	}

	return nomadClient, nil
}

// planApiJob creates a nomad api client, and runs a plan
// for the given job, printing the output diff.
func planApiJob(nomadClient *nomadApi.Client, job *nomadApi.Job) error {
	planResp, _, nomadErr := nomadClient.Jobs().Plan(job, true, nil)
	if nomadErr != nil {
		log.Errorf("Error planning job: %s", nomadErr)
		return fmt.Errorf(fmt.Sprintf("Error planning job: %s", nomadErr))
	}

	desired := planResp.Annotations.DesiredTGUpdates[*job.Name]
	logPayload := fmt.Sprintf("%+v", desired)
	if desired.Ignore > 0 {
		logPayload = "IGNORE"
	}

	log.Infof("Sucessfully planned nomad job %s - %s\n", *job.Name, logPayload)
	return nil
}

// submitApiJob creates a nomad api client, and submits the job to nomad.
func submitApiJob(nomadClient *nomadApi.Client, job *nomadApi.Job) error {
	_, _, err := nomadClient.Jobs().Register(job, nil)
	if err != nil {
		return oops.Wrapf(err, "error submitting job: %+v", job)
	}

	log.Infof("Sucessfully submitted nomad job %s\n", *job.Name)
	return nil
}

func SubmitApiJobSpecial(job *nomadApi.Job) error {
	cl, err := newNomadClient()
	if err != nil {
		return oops.Wrapf(err, "error creating nomad api client")
	}

	if err = submitApiJob(cl, job); err != nil {
		return oops.Wrapf(err, "error submitting api job")
	}

	return nil
}

func RunDeploy(j *Job, confirm, force, verbose bool) error {
	cl, err := newNomadClient()
	if err != nil {
		return oops.Wrapf(err, "error creating nomad api client")
	}

	_, aj, err := j.ToNomadJob(force)
	if err != nil {
		return oops.Wrapf(err, "error creating api job for job: %+v", j)
	}

	if err = planApiJob(cl, aj); err != nil {
		return oops.Wrapf(err, "error planning api job")
	}

	if confirm {
		if err = submitApiJob(cl, aj); err != nil {
			return oops.Wrapf(err, "error submitting api job")
		}
	} else {
		log.Infof("Skipping job submission")
	}

	return nil
}
