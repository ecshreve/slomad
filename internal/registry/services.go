package registry

import (
	_ "embed"
	"os"
	"strings"

	smd "github.com/ecshreve/slomad/pkg/slomad"
)

func promConfigHelper(tmpl string) string {
	s := strings.Replace(tmpl, "<CONSUL_TARGET>", os.Getenv("CONSUL_TARGET"), -1)
	return s
}

//go:embed config/promtail.yml
var promtailConfig string

var PromtailJob = smd.Job{
	Name:  "promtail",
	Type:  smd.SYSTEM,
	Ports: smd.BasicPortConfig(3200),
	Shape: smd.TINY_TASK,
	Env:   map[string]string{"HOSTNAME": "${attr.unique.hostname}"},
	Args: []string{
		"-config.file=/local/config/promtail.yml",
		"-server.http-listen-port=${NOMAD_PORT_http}",
	},
	Templates: map[string]string{"promtail.yml": promConfigHelper(promtailConfig)},
	Volumes: []smd.Volume{
		{Src: "/opt/nomad/data/", Dst: "/nomad/"},
		{Src: "/data/promtail", Dst: "/data"},
	},
}

var NodeExporterJob = smd.Job{
	Name:  "node-exporter",
	Type:  smd.SYSTEM,
	Ports: smd.BasicPortConfig(9100),
	Shape: smd.TINY_TASK,
	Args: []string{
		"--web.listen-address=:${NOMAD_PORT_http}",
		"--path.procfs=/host/proc",
		"--path.sysfs=/host/sys",
		"--collector.filesystem.ignored-mount-points",
		"^/(sys|proc|dev|host|etc|rootfs/var/lib/docker/containers|rootfs/var/lib/docker/overlay2|rootfs/run/docker/netns|rootfs/var/lib/docker/aufs)($$|/)",
	},
	Volumes: []smd.Volume{
		{Src: "/proc", Dst: "/host/proc"},
		{Src: "/sys", Dst: "/host/sys"},
		{Src: "/", Dst: "/rootfs"},
	},
}

var GrafanaJob = smd.Job{
	Name:    "grafana",
	Type:    smd.SERVICE,
	Target:  smd.WORKER,
	Ports:   smd.BasicPortConfig(3000),
	Shape:   smd.LARGE_TASK,
	User:    "root",
	Env:     map[string]string{"GF_SERVER_HTTP_PORT": "${NOMAD_PORT_http}"},
	Volumes: []smd.Volume{{Src: "grafana-vol", Dst: "/var/lib/grafana", Mount: true}},
}

var LokiJob = smd.Job{
	Name:   "loki",
	Type:   smd.SERVICE,
	Target: smd.WORKER,
	Ports:  smd.BasicPortConfig(3100),
	Shape:  smd.TINY_TASK,
}

//go:embed config/prometheus.yml
var prometheusConfig string

// PrometheusJob is a Job for the Prometheus service.
var PrometheusJob = smd.Job{
	Name:      "prometheus",
	Type:      smd.SERVICE,
	Target:    smd.WORKER,
	Ports:     smd.BasicPortConfig(9090),
	Shape:     smd.LARGE_TASK,
	Templates: map[string]string{"prometheus.yml": promConfigHelper(prometheusConfig)},
	Volumes:   []smd.Volume{{Src: "local/config", Dst: "/etc/prometheus"}},
}

var SpeedtestJob = smd.Job{
	Name:   "speedtest",
	Type:   smd.SERVICE,
	Target: smd.WORKER,
	Ports:  smd.BasicPortConfig(80),
	Shape:  smd.TINY_TASK,
}

var WhoamiJob = smd.Job{
	Name:   "whoami",
	Type:   smd.SERVICE,
	Target: smd.WORKER,
	Shape:  smd.TINY_TASK,
	Args:   []string{"--port", "${NOMAD_PORT_http}"},
	Ports:  smd.BasicPortConfig(80),
}

// TODO: mount nomad volume and persist data
var InfluxDBJob = smd.Job{
	Name:    "influxdb",
	Type:    smd.SERVICE,
	Target:  smd.WORKER,
	Ports:   smd.BasicPortConfig(8086),
	Shape:   smd.LARGE_TASK,
	Volumes: []smd.Volume{{Src: "influx_data", Dst: "/var/lib/influxdb"}},
}
