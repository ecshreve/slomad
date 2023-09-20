package registry

import (
	"github.com/ecshreve/slomad/pkg/slomad"
)

// TraefikJob is the job definition for the traefik service.
//
// @DEPRECATED: This job is no longer used.
var TraefikJob = slomad.Job{
	Name:   "traefik",
	Type:   slomad.SERVICE,
	Target: slomad.ALL,
	Shape:  slomad.DEFAULT_TASK,
	Ports: []*slomad.Port{
		{Label: "web", To: 80, From: 81, Static: true},
		{Label: "admin", To: 8080, From: 8081, Static: true},
	},
	Args: []string{
		"--entryPoints.web.address=:80",
		"--entryPoints.admin.address=:8081",
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
	GroupServices: map[string]string{"traefik-web": "web"},
	TaskServiceTags: map[string][]string{"traefik": {
		"traefik.enable=true",
		"traefik.http.routers.api.rule=Host(`traefik.slab.lan`)",
		"traefik.http.routers.api.service=api@internal",
	}},
}
