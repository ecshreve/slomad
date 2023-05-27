package main

import (
	"github.com/ecshreve/slomad/pkg/registry"
	"github.com/ecshreve/slomad/pkg/slomad"
	log "github.com/sirupsen/logrus"
)

func main() {
	// if err := registry.DeployTraefikJob(); err != nil {
	// 	log.Panic(err)
	// }

	// if err := registry.DeployGlances(); err != nil {
	// 	log.Panic(err)
	// }

	services := []*slomad.Job{
		// registry.GrafanaJob,
		// registry.PrometheusJob,
		// registry.LokiJob,
		// registry.ProxmoxExporterJob,
		// registry.PromtailJob,
		// registry.NodeExporterJob,
		// registry.WhoamiJob,
		// registry.SpeedtestJob,
		// registry.ControllerJob,
		// registry.NodeJob,
		// registry.JenkinsJob,
		// registry.SemaphoreJob,
		// registry.GlancesJob,
		registry.InfluxJob,
	}

	for _, x := range services {
		// if err := registry.CreateVolume(x.Name); err != nil {
		// 	log.Panic(err)
		// }

		if err := x.Plan(false); err != nil {
			log.Panic(err)
		}

		// if err := x.Deploy(false); err != nil {
		// 	log.Panic(err)
		// }
	}
	log.Info("done!")
}
