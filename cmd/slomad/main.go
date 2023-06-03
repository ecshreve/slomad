package main

import (
	"github.com/ecshreve/slomad/pkg/registry"
	"github.com/ecshreve/slomad/pkg/slomad"
	"github.com/samsarahq/go/oops"
	log "github.com/sirupsen/logrus"
)

func main() {

	services := []*slomad.App{
		registry.LokiJob,
		registry.WhoamiJob,
		registry.SpeedtestJob,
		registry.GrafanaJob,
		// 		registry.PrometheusJob,
		// 		registry.ProxmoxExporterJob,
		// 		registry.PromtailJob,
		// 		registry.NodeExporterJob,
		// 		registry.ControllerJob,
		// 		registry.NodeJob,
		// 		registry.JenkinsJob,
		// 		registry.SemaphoreJob,
		// 		registry.GlancesJob,
		// 		registry.InfluxJob,
	}

	for _, srvc := range services {
		if err := srvc.PlanApp(false); err != nil {
			log.Fatalln(oops.Wrapf(err, "error planning api job"))
		}
	}
	//something

	// 	if err := registry.DeployTraefikJob(); err != nil {
	// 		log.Panic(err)
	// 	}

	// 	if err := registry.DeployGlances(); err != nil {
	// 		log.Panic(err)
	// 	}

	// 	for _, x := range services {
	// 		if err := registry.CreateVolume(x.Name); err != nil {
	// 			log.Panic(err)
	// 		}

	// 		if err := x.Plan(false); err != nil {
	// 			log.Panic(err)
	// 		}

	//		if err := x.Deploy(false); err != nil {
	//			log.Panic(err)
	//		}
	//	}
	//
	// log.Info("done!")
}
