package registry_test

import (
	"testing"

	"github.com/ecshreve/slomad/internal/registry"
	"github.com/ecshreve/slomad/pkg/slomad"
	"github.com/samsarahq/go/snapshotter"
)

func TestRegistry(t *testing.T) {
	snap := snapshotter.New(t)
	defer snap.Verify()

	allJobs := []slomad.Job{
		registry.LokiJob,
		registry.WhoamiJob,
		registry.SpeedtestJob,
		registry.GrafanaJob,
		registry.NodeExporterJob,
		// registry.InfluxDBJob,
		// registry.PrometheusJob,
		// registry.NodeJob,
		// registry.PromtailJob,
		// registry.ControllerJob,
	}

	for _, job := range allJobs {
		snap.Snapshot(job.Name, job)
	}
}
