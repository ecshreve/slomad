package registry

import (
	_ "embed"

	smd "github.com/ecshreve/slomad/pkg/slomad"
)

//go:embed config/promtail.yml
var promtailConfig string

var PromtailJob = smd.NewAppJob(smd.JobParams{
	Name: "promtail",
	Type: smd.SYSTEM,
	TaskConfigParams: smd.TaskConfigParams{
		Ports: smd.BasicPortConfig(3200),
		Shape: smd.TINY_TASK,
		Env:   map[string]string{"HOSTNAME": "${attr.unique.hostname}"},
		Args: []string{
			"-config.file=/local/config/promtail.yml",
			"-server.http-listen-port=${NOMAD_PORT_http}",
		},
		Templates: map[string]string{"promtail.yml": promtailConfig},
	},
	StorageParams: smd.StorageParams{
		Volumes: []smd.Volume{
			{Src: "/opt/nomad/data/", Dst: "/nomad/"},
			{Src: "/data/promtail", Dst: "/data"},
		},
	},
})

var NodeExporterJob = smd.NewAppJob(smd.JobParams{
	Name: "node-exporter",
	Type: smd.SYSTEM,
	TaskConfigParams: smd.TaskConfigParams{
		Ports: smd.BasicPortConfig(9100),
		Shape: smd.TINY_TASK,
		Args: []string{
			"--web.listen-address=:${NOMAD_PORT_http}",
			"--path.procfs=/host/proc",
			"--path.sysfs=/host/sys",
			"--collector.filesystem.ignored-mount-points",
			"^/(sys|proc|dev|host|etc|rootfs/var/lib/docker/containers|rootfs/var/lib/docker/overlay2|rootfs/run/docker/netns|rootfs/var/lib/docker/aufs)($$|/)",
		},
	},
	StorageParams: smd.StorageParams{
		Volumes: []smd.Volume{
			{Src: "/proc", Dst: "/host/proc"},
			{Src: "/sys", Dst: "/host/sys"},
			{Src: "/", Dst: "/rootfs"},
		},
	},
})

var GrafanaJob = smd.NewAppJob(smd.JobParams{
	Name:   "grafana",
	Type:   smd.SERVICE,
	Target: smd.WORKER,
	TaskConfigParams: smd.TaskConfigParams{
		Ports: smd.BasicPortConfig(3000),
		Shape: smd.LARGE_TASK,
		User:  "root",
		Env:   map[string]string{"GF_SERVER_HTTP_PORT": "${NOMAD_PORT_http}"},
	},
	StorageParams: smd.StorageParams{
		Storage: smd.StringPtr("grafana"),
		Volumes: []smd.Volume{{Src: "grafana-vol", Dst: "/var/lib/grafana", Mount: true}},
	},
})

var LokiJob = smd.NewAppJob(smd.JobParams{
	Name:   "loki",
	Type:   smd.SERVICE,
	Target: smd.WORKER,
	TaskConfigParams: smd.TaskConfigParams{
		Ports: smd.BasicPortConfig(3100),
		Shape: smd.TINY_TASK,
	},
})

//go:embed config/prometheus.yml
var prometheusConfig string

var PrometheusJob = smd.NewAppJob(smd.JobParams{
	Name:   "prometheus",
	Type:   smd.SERVICE,
	Target: smd.WORKER,
	TaskConfigParams: smd.TaskConfigParams{
		Ports:     smd.BasicPortConfig(9090),
		Shape:     smd.LARGE_TASK,
		Templates: map[string]string{"prometheus.yml": prometheusConfig},
	},
	StorageParams: smd.StorageParams{
		Storage: smd.StringPtr("prometheus"),
		Volumes: []smd.Volume{{Src: "local/config", Dst: "/etc/prometheus"}},
	},
})

var SpeedtestJob = smd.NewAppJob(smd.JobParams{
	Name:   "speedtest",
	Type:   smd.SERVICE,
	Target: smd.WORKER,
	TaskConfigParams: smd.TaskConfigParams{
		Ports: smd.BasicPortConfig(80),
		Shape: smd.TINY_TASK,
	},
})

var WhoamiJob = smd.NewAppJob(smd.JobParams{
	Name:   "whoami",
	Type:   smd.SERVICE,
	Target: smd.WORKER,
	TaskConfigParams: smd.TaskConfigParams{
		Shape: smd.TINY_TASK,
		Args:  []string{"--port", "${NOMAD_PORT_http}"},
		Ports: smd.BasicPortConfig(80),
	},
})

// TODO: mount nomad volume and persist data
var InfluxDBJob = smd.NewAppJob(smd.JobParams{
	Name:   "influxdb",
	Type:   smd.SERVICE,
	Target: smd.WORKER,
	TaskConfigParams: smd.TaskConfigParams{
		Ports: smd.BasicPortConfig(8086),
		Shape: smd.LARGE_TASK,
	},
	StorageParams: smd.StorageParams{
		Volumes: []smd.Volume{
			{Src: "influx_data", Dst: "/var/lib/influxdb"},
		},
	},
})

var ControllerJob = smd.NewStorageJob(smd.JobParams{
	Name:   "storage-controller",
	Type:   smd.SERVICE,
	Target: smd.WORKER,
	TaskConfigParams: smd.TaskConfigParams{
		Ports: smd.BasicPortConfig(0),
		Shape: smd.DEFAULT_TASK,
		Args:  getStorageArgs("controller"),
	},
})

var NodeJob = smd.NewStorageJob(smd.JobParams{
	Name: "storage-node",
	Type: smd.SYSTEM,
	TaskConfigParams: smd.TaskConfigParams{
		Ports: smd.BasicPortConfig(0),
		Shape: smd.TINY_TASK,
		Args:  getStorageArgs("node"),
	},
})