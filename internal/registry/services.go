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
	Shape: smd.XTINY_TASK,
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
	Shape: smd.XTINY_TASK,
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

// TODO: mount nomad volume and persist data
var GrafanaJob = smd.Job{
	Name:   "grafana",
	Type:   smd.SERVICE,
	Target: smd.NODE,
	Ports:  smd.BasicPortConfig(3000),
	Shape:  smd.LARGE_TASK,
	User:   "root",
	Env: map[string]string{
		"GF_SERVER_ROOT_URL":            "http://slab.lan/grafana",
		"GF_SERVER_SERVE_FROM_SUB_PATH": "true",
	},
	// Volumes: []smd.Volume{{Src: "grafana-vol", Dst: "/var/lib/grafana", Mount: true}},
}

var LokiJob = smd.Job{
	Name:   "loki",
	Type:   smd.SERVICE,
	Target: smd.NODE,
	Ports: []*smd.Port{
		{Label: "http", To: 3100, From: 3100, Static: true},
	},
	Shape: smd.TINY_TASK,
}

//go:embed config/prometheus.yml
var prometheusConfig string

// PrometheusJob is a Job for the Prometheus service.
var PrometheusJob = smd.Job{
	Name:      "prometheus",
	Type:      smd.SERVICE,
	Target:    smd.NODE,
	Ports:     smd.BasicPortConfig(9090),
	Shape:     smd.LARGE_TASK,
	Templates: map[string]string{"prometheus.yml": promConfigHelper(prometheusConfig)},
	Volumes:   []smd.Volume{{Src: "local/config", Dst: "/etc/prometheus"}},
	Args: []string{
		"--web.external-url=http://slab.lan/prometheus",
		"--config.file=/etc/prometheus/prometheus.yml",
	},
}

var SpeedtestJob = smd.Job{
	Name:   "speedtest",
	Type:   smd.SERVICE,
	Target: smd.NODE,
	Ports:  smd.BasicPortConfig(80),
	Shape:  smd.XTINY_TASK,
}

var UptimeJob = smd.Job{
	Name:   "uptime",
	Type:   smd.SERVICE,
	Target: smd.NODE,
	Ports: []*smd.Port{
		{Label: "http", To: 3001, From: 3001, Static: true},
	},
	Shape: smd.DEFAULT_TASK,
}

var WhoamiJob = smd.Job{
	Name:   "whoami",
	Type:   smd.SERVICE,
	Target: smd.NODE,
	Shape:  smd.XXTINY_TASK,
	Ports:  smd.BasicPortConfig(80),
	TaskServiceTags: map[string][]string{
		"whoami": {
			"urlprefix-/whoami",
		},
	},
}

// // TODO: mount nomad volume and persist data
// var InfluxDBJob = smd.Job{
// 	Name:    "influxdb",
// 	Type:    smd.SERVICE,
// 	Target:  smd.WORKER,
// 	Ports:   smd.BasicPortConfig(8086),
// 	Shape:   smd.LARGE_TASK,
// 	Volumes: []smd.Volume{{Src: "influx_data", Dst: "/var/lib/influxdb"}},
// }
