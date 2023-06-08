package main

import (
	"fmt"
	"os"

	"github.com/ecshreve/slomad/pkg/slomad"
	nomadApi "github.com/hashicorp/nomad/api"
	"github.com/kylelemons/godebug/pretty"
	"github.com/samsarahq/go/oops"
	log "github.com/sirupsen/logrus"
)

// RunDeploy runs a deploy for the given job.
func RunDeploy(j *slomad.Job, confirm, force, verbose bool) error {
	aj, err := j.GetNomadApiJob(force)
	if err != nil {
		return oops.Wrapf(err, "error creating api job for slomad job: %+v", j)
	}

	cl, err := newNomadClient()
	if err != nil {
		return oops.Wrapf(err, "error creating nomad api client")
	}

	ignore, err := planApiJob(cl, aj, verbose)
	if err != nil {
		return oops.Wrapf(err, "error planning api job")
	}

	if ignore {
		log.Infof("skipping submit for job: %s", *aj.Name)
		return nil
	}

	if confirm {
		if err = submitApiJob(cl, aj); err != nil {
			return oops.Wrapf(err, "error submitting api job")
		}
	}

	return nil
}

// newNomadClient returns a default nomad api client.
func newNomadClient() (*nomadApi.Client, error) {
	nomadConfig := nomadApi.DefaultConfig()
	nomadConfig.Address = os.Getenv("NOMAD_TARGET")
	nomadClient, err := nomadApi.NewClient(nomadConfig)
	if err != nil {
		return nil, oops.Wrapf(err, "unable to create nomad api client")
	}

	return nomadClient, nil
}

// planApiJob runs a plan for the given job, and returns whether or not the job
// can be ignored, and any errors encountered.
func planApiJob(nomadClient *nomadApi.Client, job *nomadApi.Job, diff bool) (bool, error) {
	planResp, _, nomadErr := nomadClient.Jobs().Plan(job, true, nil)
	if nomadErr != nil {
		log.Errorf("Error planning job: %s", nomadErr)
		return false, fmt.Errorf(fmt.Sprintf("Error planning job: %s", nomadErr))
	}

	desired := planResp.Annotations.DesiredTGUpdates[*job.Name]
	logPayload := fmt.Sprintf("%+v", desired)
	if desired.Ignore > 0 {
		logPayload = "IGNORE"
	} else if diff {
		log.Infof("Plan diff for nomad job %s:\n", *job.Name)
		pretty.Print(planResp.Diff)
	}

	log.Infof("Successfully planned nomad job %s - %s\n", *job.Name, logPayload)

	return desired.Ignore > 0, nil
}

// submitApiJob submits the given job to the nomad cluster.
func submitApiJob(nomadClient *nomadApi.Client, job *nomadApi.Job) error {
	_, _, err := nomadClient.Jobs().Register(job, nil)
	if err != nil {
		return oops.Wrapf(err, "error submitting job: %+v", job)
	}

	log.Infof("Successfully submitted nomad job %s\n", *job.Name)
	return nil
}
