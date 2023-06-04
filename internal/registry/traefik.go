package registry

import (
	"fmt"
	"time"

	"github.com/ecshreve/slomad/pkg/slomad"
	nomadStructs "github.com/hashicorp/nomad/nomad/structs"
	"github.com/samsarahq/go/oops"
	log "github.com/sirupsen/logrus"
)

var TraefikJob = nomadStructs.Job{
	ID:          "traefik",
	Name:        "traefik",
	Region:      "global",
	Priority:    92,
	Datacenters: []string{"dcs"},
	Type:        "service",
	TaskGroups: []*nomadStructs.TaskGroup{
		{
			Services: []*nomadStructs.Service{
				{
					Name:      "traefik-web",
					PortLabel: "web",
					Checks: []*nomadStructs.ServiceCheck{
						{
							Name:          fmt.Sprintf("%s = tcp check", "traefik-web"),
							Type:          nomadStructs.ServiceCheckTCP,
							Interval:      10 * time.Second,
							Timeout:       2 * time.Second,
							InitialStatus: "passing",
						},
					},
					Provider: "consul",
				},
			},
			Name:  "traefik",
			Count: 1,
			Tasks: []*nomadStructs.Task{
				{
					Name:   "traefik",
					Driver: "docker",
					Config: map[string]interface{}{
						"image":        "reg.slab.lan:5000/traefik:latest",
						"network_mode": "host",
						"args": []string{
							"--entryPoints.web.address=:80",
							"--entryPoints.websecure.address=:443",
							"--entryPoints.admin.address=:8081",
							"--entrypoints.websecure.http.redirections.entryPoint.to=web",
							"--entrypoints.websecure.http.redirections.entryPoint.scheme=http",
							"--accesslog=true",
							"--api=true",
							"--api.dashboard=true",
							"--api.insecure=true",
							"--ping=true",
							"--ping.entryPoint=admin",
							"--providers.consulcatalog=true",
							"--providers.consulcatalog.endpoint.address=10.35.220.50:8500",
							"--providers.consulcatalog.prefix=traefik",
							"--providers.consulcatalog.refreshInterval=30s",
							"--providers.consulcatalog.exposedByDefault=false",
							"--providers.consulcatalog.defaultrule=Host(`{{ .Name }}.slabstaging.lan`)",
							"--providers.consulcatalog.endpoint.tls.insecureskipverify=true",
						},
					},
					Resources: &nomadStructs.Resources{
						CPU:      512,
						MemoryMB: 512,
					},
					LogConfig: nomadStructs.DefaultLogConfig(),
					Services: []*nomadStructs.Service{
						{
							Name:      "traefik",
							PortLabel: "websecure",
							Tags: []string{
								"traefik.enable=true",
								"traefik.http.routers.api.rule=Host(`traefik.slabstaging.lan`)",
								"traefik.http.routers.api.service=api@internal",
							},
							TaskName: "traefik",
							Checks: []*nomadStructs.ServiceCheck{
								{
									Name:          fmt.Sprintf("%s = http check", "traefik"),
									Type:          nomadStructs.ServiceCheckHTTP,
									Interval:      10 * time.Second,
									Timeout:       2 * time.Second,
									InitialStatus: "passing",
									Path:          "/ping",
									PortLabel:     "admin",
									TaskName:      "traefik",
								},
							},
							Provider: "consul",
						},
					},
				},
			},
			RestartPolicy:    nomadStructs.NewRestartPolicy("service"),
			ReschedulePolicy: &nomadStructs.DefaultServiceJobReschedulePolicy,
			EphemeralDisk: &nomadStructs.EphemeralDisk{
				SizeMB: 256,
			},
			Networks: []*nomadStructs.NetworkResource{
				{
					Mode: "host",
					ReservedPorts: []nomadStructs.Port{
						{
							Label: "web",
							Value: 80,
							To:    0,
						},
						{
							Label: "websecure",
							Value: 443,
							To:    0,
						},
						{
							Label: "admin",
							Value: 8081,
							To:    0,
						},
					},
				},
			},
		},
	},
	Namespace: "default",
	Constraints: []*nomadStructs.Constraint{
		{
			LTarget: "${attr.unique.hostname}",
			RTarget: "worker-0",
			Operand: "regexp",
		},
	},
}

// DeployTraefikJob deploys the Traefik job to Nomad.
func DeployTraefikJob(confirm bool) error {
	job := &TraefikJob
	if err := job.Validate(); err != nil {
		log.Errorf("Nomad job validation failed. Error: %s\n", err)
		return err
	}

	apiJob, err := slomad.ConvertJob(job)
	if err != nil {
		log.Errorf("Failed to convert nomad job in api call. Error: %s\n", err)
		return err
	}

	if confirm {
		log.Infof("deploying %s", job.Name)
		if err = slomad.SubmitApiJob(apiJob); err != nil {
			return oops.Wrapf(err, "error submitting api job")
		}
	} else {
		log.Debugf("skipping deploy %s", job.Name)
	}

	return nil
}
