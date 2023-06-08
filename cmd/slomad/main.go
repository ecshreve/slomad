package main

import (
	"os"

	"github.com/ecshreve/slomad/internal/registry"
	"github.com/ecshreve/slomad/pkg/slomad"
	"github.com/samsarahq/go/oops"
	log "github.com/sirupsen/logrus"
)

func main() {
	args := os.Args[1:]
	confirm := false
	diff := false
	if len(args) > 0 {
		switch args[0] {
		case "confirm":
			confirm = true
		case "diff":
			diff = true
		default:
			log.Warnf("unknown arg: %s", args[0])
		}
	}

	// tt := slomad.NewTTJob()
	// ttApi, _ := tt.GetNomadApiJob(false)
	// exp, _ := registry.GetTraefikJob()

	// diffs := pretty.Diff(ttApi, exp)
	// for _, d := range diffs {
	// 	fmt.Println(d)
	// }

	// if err := RunDeploy(tt, true, false, false); err != nil {
	// 	log.Fatal(err)
	// }

	// if err := RunTraefikDeploy(confirm); err != nil {
	// 	log.Fatalln(oops.Wrapf(err, "error deploying traefik job"))
	// }

	services := []slomad.Job{
		// registry.LokiJob,
		// registry.WhoamiJob,
		// registry.SpeedtestJob,
		// registry.GrafanaJob,
		// registry.PrometheusJob,
		// registry.NodeExporterJob,
		// registry.PromtailJob,
		// registry.ControllerJob,
		// registry.NodeJob,
		// registry.InfluxDBJob,
		// registry.PlexJob,
		registry.TraefikJob,
	}

	for _, srvc := range services {
		if err := RunDeploy(&srvc, confirm, false, diff); err != nil {
			log.Fatalln(oops.Wrapf(err, "error deploying job"))
		}
	}
}
