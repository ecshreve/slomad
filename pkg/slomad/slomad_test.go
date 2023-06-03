package slomad_test

import (
	"testing"

	"github.com/ecshreve/slomad/pkg/slomad"
	"github.com/samsarahq/go/snapshotter"
)

func TestSlomad(t *testing.T) {
	snap := snapshotter.New(t)
	defer snap.Verify()

	testJob := slomad.NewAppJob(slomad.JobParams{
		Name:   "test-job",
		Type:   slomad.SERVICE,
		Target: slomad.WORKER,
		TaskConfigParams: slomad.TaskConfigParams{
			Args:  []string{"echo hello"},
			Ports: []*slomad.Port{{Label: "http", To: 8080, From: 0, Static: false}},
			Shape: slomad.DEFAULT_TASK,
		},
	})

	snap.Snapshot("test-job", testJob)

	nj, aj, _ := testJob.ToNomadJob(false)
	snap.Snapshot("nomad-job", nj)
	snap.Snapshot("api-job", aj)
}
