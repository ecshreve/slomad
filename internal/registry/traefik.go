package registry

import (
	"github.com/ecshreve/slomad/pkg/slomad"
)

var TraefikJob = slomad.Job{
	Name:   "traefik",
	Type:   slomad.SERVICE,
	Target: slomad.WORKER_0,
	Shape:  slomad.DEFAULT_TASK,
	Ports: []*slomad.Port{
		{Label: "web", To: 0, From: 80, Static: true},
		{Label: "websecure", To: 0, From: 443, Static: true},
		{Label: "admin", To: 0, From: 8081, Static: true},
	},
	Args: []string{
		"--entryPoints.web.address=:80",
		"--entryPoints.websecure.address=:443",
		"--entryPoints.admin.address=:8081",
		"--entrypoints.websecure.http.redirections.entryPoint.to=web",
		"--entrypoints.websecure.http.redirections.entryPoint.scheme=http",
		"--accesslog=true",
		"--api=true",
		"--api.dashboard=true",
		"--api.insecure=true",
		"--metrics=true",
		"--metrics.prometheus=true",
		"--metrics.prometheus.addEntryPointsLabels=true",
		"--ping=true",
		"--ping.entryPoint=admin",
		"--providers.consulcatalog=true",
		"--providers.consulcatalog.endpoint.address=127.0.0.1:8500",
		"--providers.consulcatalog.prefix=traefik",
		"--providers.consulcatalog.refreshInterval=30s",
		"--providers.consulcatalog.exposedByDefault=false",
		"--providers.consulcatalog.defaultrule=Host(`{{ .Name }}.slab.lan`)",
		"--providers.consulcatalog.endpoint.tls.insecureskipverify=true",
	},
}

// var traefikJob = nomadStructs.Job{
// 	ID:          "traefik",
// 	Name:        "traefik",
// 	Region:      "global",
// 	Priority:    92,
// 	Datacenters: []string{"dcs"},
// 	Type:        "service",
// 	TaskGroups: []*nomadStructs.TaskGroup{
// 		{
// 			Services: []*nomadStructs.Service{
// 				{
// 					Name:      "traefik-web",
// 					PortLabel: "web",
// 					Checks: []*nomadStructs.ServiceCheck{
// 						{
// 							Name:          fmt.Sprintf("%s = tcp check", "traefik-web"),
// 							Type:          nomadStructs.ServiceCheckTCP,
// 							Interval:      10 * time.Second,
// 							Timeout:       2 * time.Second,
// 							InitialStatus: "passing",
// 						},
// 					},
// 					Provider: "consul",
// 				},
// 			},
// 			Name:  "traefik",
// 			Count: 1,
// 			Tasks: []*nomadStructs.Task{
// 				{
// 					Name:   "traefik",
// 					Driver: "docker",
// 					Config: map[string]interface{}{
// 						"image":        "reg.slab.lan:5000/traefik:latest",
// 						"network_mode": "host",
// 						"args": []string{
// 							"--entryPoints.web.address=:80",
// 							"--entryPoints.websecure.address=:443",
// 							"--entryPoints.admin.address=:8081",
// 							"--entrypoints.websecure.http.redirections.entryPoint.to=web",
// 							"--entrypoints.websecure.http.redirections.entryPoint.scheme=http",
// 							"--accesslog=true",
// 							"--api=true",
// 							"--api.dashboard=true",
// 							"--api.insecure=true",
// 							"--metrics=true",
// 							"--metrics.prometheus=true",
// 							"--metrics.prometheus.addEntryPointsLabels=true",
// 							"--ping=true",
// 							"--ping.entryPoint=admin",
// 							"--providers.consulcatalog=true",
// 							"--providers.consulcatalog.endpoint.address=127.0.0.1:8500",
// 							"--providers.consulcatalog.prefix=traefik",
// 							"--providers.consulcatalog.refreshInterval=30s",
// 							"--providers.consulcatalog.exposedByDefault=false",
// 							"--providers.consulcatalog.defaultrule=Host(`{{ .Name }}.slab.lan`)",
// 							"--providers.consulcatalog.endpoint.tls.insecureskipverify=true",
// 						},
// 					},
// 					Resources: &nomadStructs.Resources{
// 						CPU:      512,
// 						MemoryMB: 512,
// 					},
// 					LogConfig: nomadStructs.DefaultLogConfig(),
// 					Services: []*nomadStructs.Service{
// 						{
// 							Name:      "traefik",
// 							PortLabel: "websecure",
// 							Tags: []string{
// 								"traefik.enable=true",
// 								"traefik.http.routers.api.rule=Host(`traefik.slab.lan`)",
// 								"traefik.http.routers.api.service=api@internal",
// 							},
// 							TaskName: "traefik",
// 							Checks: []*nomadStructs.ServiceCheck{
// 								{
// 									Name:          fmt.Sprintf("%s = http check", "traefik"),
// 									Type:          nomadStructs.ServiceCheckHTTP,
// 									Interval:      10 * time.Second,
// 									Timeout:       2 * time.Second,
// 									InitialStatus: "passing",
// 									Path:          "/ping",
// 									PortLabel:     "admin",
// 									TaskName:      "traefik",
// 								},
// 							},
// 							Provider: "consul",
// 						},
// 					},
// 				},
// 			},
// 			RestartPolicy:    nomadStructs.NewRestartPolicy("service"),
// 			ReschedulePolicy: &nomadStructs.DefaultServiceJobReschedulePolicy,
// 			EphemeralDisk: &nomadStructs.EphemeralDisk{
// 				SizeMB: 256,
// 			},
// 			Networks: []*nomadStructs.NetworkResource{
// 				{
// 					Mode: "host",
// 					ReservedPorts: []nomadStructs.Port{
// 						{
// 							Label: "web",
// 							Value: 80,
// 							To:    0,
// 						},
// 						{
// 							Label: "websecure",
// 							Value: 443,
// 							To:    0,
// 						},
// 						{
// 							Label: "admin",
// 							Value: 8081,
// 							To:    0,
// 						},
// 					},
// 				},
// 			},
// 		},
// 	},
// 	Namespace: "default",
// 	Constraints: []*nomadStructs.Constraint{
// 		{
// 			LTarget: "${attr.unique.hostname}",
// 			RTarget: "worker-0",
// 			Operand: "regexp",
// 		},
// 	},
// }

// convertJob converts a Nomad Job to a Nomad API Job.
// func convertJob(in *nomadStructs.Job) (*nomadApi.Job, error) {
// 	gob.Register([]map[string]interface{}{})
// 	gob.Register([]interface{}{})

// 	var apiJob *nomadApi.Job
// 	buf := new(bytes.Buffer)
// 	if err := gob.NewEncoder(buf).Encode(in); err != nil {
// 		return nil, err
// 	}
// 	if err := gob.NewDecoder(buf).Decode(&apiJob); err != nil {
// 		return nil, err
// 	}

// 	return apiJob, nil
// }

// func GetTraefikJob() (*nomadApi.Job, error) {
// 	if err := traefikJob.Validate(); err != nil {
// 		return nil, err
// 	}
// 	return convertJob(&traefikJob)
// }
