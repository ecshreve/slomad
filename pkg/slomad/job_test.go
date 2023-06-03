package slomad_test

import (
	"testing"

	"github.com/ecshreve/slomad/pkg/slomad"
	"github.com/samsarahq/go/snapshotter"
)

func TestNewServiceJob(t *testing.T) {
	snapshotter := snapshotter.New(t)
	defer snapshotter.Verify()

	p := slomad.JobParams{
		Name:   "whoami",
		Target: slomad.WORKER,
		TaskConfigParams: slomad.TaskConfigParams{
			Ports: []slomad.PortParams{{Label: "http", To: 80}},
			Args:  []string{"--port", "${NOMAD_PORT_http}"},
			Shape: slomad.TINY_TASK,
		},
	}
	snapshotter.Snapshot("jobParams", p)

	jjj := slomad.NewServiceJob(p)
	snapshotter.Snapshot("JJJob", jjj)

	nj, aj, err := jjj.ToNomadJob()
	snapshotter.Snapshot("NomadJob", nj)
	snapshotter.Snapshot("NomadApiJob", aj)
	snapshotter.Snapshot("error", err)
}
