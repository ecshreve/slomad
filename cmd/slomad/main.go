package main

import (
	"github.com/ecshreve/slomad/pkg/slomad"
	"github.com/samsarahq/go/oops"
	log "github.com/sirupsen/logrus"
)

func main() {
	// var WhoamiJob = &slomad.Job{
	// 	Name:       "whoami",
	// 	Image:      getDockerImageString("whoami"),
	// 	Args:       []string{"--port", "${NOMAD_PORT_http}"},
	// 	CommonArgs: getCommonJobArgs("docker", "^worker-[0-9]+$", 1, 50),
	// 	Ports:      []slomad.Port{{Label: "http"}},
	// 	Size:       map[string]int{"cpu": 128, "mem": 128},
	// }

	var WhoamiJob = slomad.NewServiceJob(slomad.JobParams{
		Name:   "whoami",
		Target: slomad.WORKER,
		TaskConfigParams: slomad.TaskConfigParams{
			Shape: slomad.TINY_TASK,
			Args:  []string{"--port", "${NOMAD_PORT_http}"},
			Ports: []slomad.PortParams{{Label: "http", To: 80}},
		},
	})

	if err := WhoamiJob.PlanAndApplyJJJ(false); err != nil {
		log.Fatalln(oops.Wrapf(err, "error planning api job"))
	}
	//something

	// 	if err := registry.DeployTraefikJob(); err != nil {
	// 		log.Panic(err)
	// 	}

	// 	if err := registry.DeployGlances(); err != nil {
	// 		log.Panic(err)
	// 	}

	// 	services := []*slomad.Job{
	// 		registry.GrafanaJob,
	// 		registry.PrometheusJob,
	// 		registry.LokiJob,
	// 		registry.ProxmoxExporterJob,
	// 		registry.PromtailJob,
	// 		registry.NodeExporterJob,
	// 		registry.WhoamiJob,
	// 		registry.SpeedtestJob,
	// 		registry.ControllerJob,
	// 		registry.NodeJob,
	// 		registry.JenkinsJob,
	// 		registry.SemaphoreJob,
	// 		registry.GlancesJob,
	// 		registry.InfluxJob,
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
