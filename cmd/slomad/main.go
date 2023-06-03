package main

import (
	"os"

	"github.com/ecshreve/slomad/pkg/registry"
	"github.com/ecshreve/slomad/pkg/slomad"
	"github.com/samsarahq/go/oops"
	log "github.com/sirupsen/logrus"
)

func main() {

	services := []*slomad.Job{
		registry.LokiJob,
		registry.WhoamiJob,
		registry.SpeedtestJob,
		registry.GrafanaJob,
		registry.PrometheusJob,
		registry.NodeExporterJob,
		registry.PromtailJob,
		registry.ControllerJob,
		registry.NodeJob,
		registry.InfluxDBJob,
	}

	args := os.Args[1:]
	apply := false
	if len(args) > 0 {
		switch args[0] {
		case "confirm":
			apply = true

		default:
			log.Warnf("unknown arg: %s", args[0])
		}
	}

	for _, srvc := range services {
		if err := srvc.Plan(false); err != nil {
			log.Fatalln(oops.Wrapf(err, "error planning api job"))
		}

		if apply {
			log.Infof("deploying %s", srvc.Name)
			if err := srvc.Deploy(false); err != nil {
				log.Fatalln(oops.Wrapf(err, "error submitting api job"))
			}
		} else {
			log.Debugf("skipping deploy %s", srvc.Name)
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
