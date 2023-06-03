package registry

import (
	_ "embed"

	"github.com/ecshreve/slomad/pkg/slomad"
)

//go:embed config/promtail.yml
var promtailConfig string

var PromtailJob = slomad.NewAppJob(slomad.JobParams{
	Name: "promtail",
	Type: slomad.SYSTEM,
	TaskConfigParams: slomad.TaskConfigParams{
		Ports: slomad.BasicPorts(3200),
		Shape: slomad.TINY_TASK,
		Env:   map[string]string{"HOSTNAME": "${attr.unique.hostname}"},
		Args: []string{
			"-config.file=/local/config/promtail.yml",
			"-server.http-listen-port=${NOMAD_PORT_http}",
		},
		Templates: map[string]string{"promtail.yml": promtailConfig},
	},
	StorageParams: slomad.StorageParams{
		Volumes: []slomad.Volume{
			{Src: "/opt/nomad/data/", Dst: "/nomad/"},
			{Src: "/data/promtail", Dst: "/data"},
		},
	},
})

var NodeExporterJob = slomad.NewAppJob(slomad.JobParams{
	Name: "node-exporter",
	Type: slomad.SYSTEM,
	TaskConfigParams: slomad.TaskConfigParams{
		Ports: slomad.BasicPorts(9100),
		Shape: slomad.TINY_TASK,
		Args: []string{
			"--web.listen-address=:${NOMAD_PORT_http}",
			"--path.procfs=/host/proc",
			"--path.sysfs=/host/sys",
			"--collector.filesystem.ignored-mount-points",
			"^/(sys|proc|dev|host|etc|rootfs/var/lib/docker/containers|rootfs/var/lib/docker/overlay2|rootfs/run/docker/netns|rootfs/var/lib/docker/aufs)($$|/)",
		},
	},
	StorageParams: slomad.StorageParams{
		Volumes: []slomad.Volume{
			{Src: "/proc", Dst: "/host/proc"},
			{Src: "/sys", Dst: "/host/sys"},
			{Src: "/", Dst: "/rootfs"},
		},
	},
})

var GrafanaJob = slomad.NewAppJob(slomad.JobParams{
	Name:   "grafana",
	Type:   slomad.SERVICE,
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

var LokiJob = slomad.NewAppJob(slomad.JobParams{
	Name:   "loki",
	Type:   slomad.SERVICE,
	Target: slomad.WORKER,
	TaskConfigParams: slomad.TaskConfigParams{
		Ports: slomad.BasicPorts(3100),
		Shape: slomad.TINY_TASK,
	},
})

//go:embed config/prometheus.yml
var prometheusConfig string

var PrometheusJob = slomad.NewAppJob(slomad.JobParams{
	Name:   "prometheus",
	Type:   slomad.SERVICE,
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

var SpeedtestJob = slomad.NewAppJob(slomad.JobParams{
	Name:   "speedtest",
	Type:   slomad.SERVICE,
	Target: slomad.WORKER,
	TaskConfigParams: slomad.TaskConfigParams{
		Ports: slomad.BasicPorts(80),
		Shape: slomad.TINY_TASK,
	},
})

var WhoamiJob = slomad.NewAppJob(slomad.JobParams{
	Name:   "whoami",
	Type:   slomad.SERVICE,
	Target: slomad.WORKER,
	TaskConfigParams: slomad.TaskConfigParams{
		Shape: slomad.TINY_TASK,
		Args:  []string{"--port", "${NOMAD_PORT_http}"},
		Ports: slomad.BasicPorts(80),
	},
})
