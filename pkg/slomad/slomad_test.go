package slomad_test

import (
	"testing"

	"github.com/ecshreve/slomad/pkg/slomad"
	"github.com/samsarahq/go/snapshotter"
	"github.com/stretchr/testify/assert"
)

func TestSlomad(t *testing.T) {
	snap := snapshotter.New(t)
	defer snap.Verify()

	testJob := slomad.NewJob(slomad.JobParams{
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

	aj, _ := testJob.GetNomadApiJob(false)
	snap.Snapshot("api-job", aj)

	aj2, _ := testJob.GetNomadApiJob(true)
	assert.NotEqual(t, aj.Meta["run_uuid"], aj2.Meta["run_uuid"])
}

func TestSlomadRegistryJob(t *testing.T) {
	snap := snapshotter.New(t)
	defer snap.Verify()

	var GrafanaJob = slomad.NewJob(slomad.JobParams{
		Name:   "grafana",
		Type:   slomad.SERVICE,
		Target: slomad.WORKER,
		TaskConfigParams: slomad.TaskConfigParams{
			Ports: slomad.BasicPortConfig(3000),
			Shape: slomad.LARGE_TASK,
			User:  "root",
			Env:   map[string]string{"GF_SERVER_HTTP_PORT": "${NOMAD_PORT_http}"},
		},
		StorageParams: slomad.StorageParams{
			Volumes: []slomad.Volume{{Src: "grafana-vol", Dst: "/var/lib/grafana", Mount: true}},
		},
	})

	snap.Snapshot("grafana-job", GrafanaJob)

	aj, _ := GrafanaJob.GetNomadApiJob(false)
	snap.Snapshot("api-job", aj)

}
