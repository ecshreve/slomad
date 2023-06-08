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
	GroupServices: map[string]string{"traefik-web": "web"},
	TaskServiceTags: map[string][]string{"traefik": {
		"traefik.enable=true",
		"traefik.http.routers.api.rule=Host(`traefik.slab.lan`)",
		"traefik.http.routers.api.service=api@internal",
	}},
}
