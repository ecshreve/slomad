package registry

import (
	_ "embed"

	"github.com/ecshreve/slomad/pkg/slomad"
)

var GrafanaJob = slomad.NewServiceJob(slomad.JobParams{
	Name:   "grafana",
	Target: slomad.WORKER,
	TaskConfigParams: slomad.TaskConfigParams{
		Ports: slomad.BasicPorts(3000),
		Shape: slomad.LARGE_TASK,
		User:  "root",
		Env:   map[string]string{"GF_SERVER_HTTP_PORT": "${NOMAD_PORT_http}"},
	},
	StorageParams: slomad.StorageParams{
		Storage: slomad.StringPtr("grafana"),
		Volumes: []slomad.Volume{{Src: "grafana-vol", Dst: "/var/lib/grafana", Mount: true}},
	},
})

var LokiJob = slomad.NewServiceJob(slomad.JobParams{
	Name:   "loki",
	Target: slomad.WORKER,
	TaskConfigParams: slomad.TaskConfigParams{
		Ports: slomad.BasicPorts(3100),
		Shape: slomad.TINY_TASK,
	},
})

//go:embed config/prometheus.yml
var prometheusConfig string

var PrometheusJob = slomad.NewServiceJob(slomad.JobParams{
	Name:   "prometheus",
	Target: slomad.WORKER,
	TaskConfigParams: slomad.TaskConfigParams{
		Ports:     slomad.BasicPorts(9090),
		Shape:     slomad.LARGE_TASK,
		Templates: map[string]string{"prometheus.yml": prometheusConfig},
	},
	StorageParams: slomad.StorageParams{
		Storage: slomad.StringPtr("prometheus"),
		Volumes: []slomad.Volume{{Src: "local/config", Dst: "/etc/prometheus"}},
	},
})

var SpeedtestJob = slomad.NewServiceJob(slomad.JobParams{
	Name:   "speedtest",
	Target: slomad.WORKER,
	TaskConfigParams: slomad.TaskConfigParams{
		Ports: slomad.BasicPorts(80),
		Shape: slomad.TINY_TASK,
	},
})

var WhoamiJob = slomad.NewServiceJob(slomad.JobParams{
	Name:   "whoami",
	Target: slomad.WORKER,
	TaskConfigParams: slomad.TaskConfigParams{
		Shape: slomad.TINY_TASK,
		Args:  []string{"--port", "${NOMAD_PORT_http}"},
		Ports: slomad.BasicPorts(80),
	},
})
