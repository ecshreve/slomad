package registry

import (
	_ "embed"

	"github.com/ecshreve/slomad/pkg/slomad"
)

var ProxmoxExporterJob = &slomad.Job{
	Name:       "pve-exporter",
	Image:      getDockerImageString("pve-exporter"),
	CommonArgs: getCommonJobArgs("docker", "^worker-[0-9]+$", 1, 50),
	Ports:      []slomad.Port{{Label: "http", To: 9221}},
	Size:       map[string]int{"cpu": 128, "mem": 128},
	Env: map[string]string{
		"PVE_USER":        "jenkins@pam",
		"PVE_TOKEN_NAME":  "jenkins@pam!jenkins",
		"PVE_TOKEN_VALUE": "5f355f17-f640-4807-bf5b-1a1aa6262506",
		"PVE_VERIFY_SSL":  "false",
	},
}

var NodeExporterJob = &slomad.Job{
	Name:       "nodeexporter",
	Image:      getDockerImageString("node_exporter"),
	JobType:    "system",
	CommonArgs: getCommonJobArgs("docker", "^*$", 1, 50),
	Ports:      []slomad.Port{{Label: "http"}},
	Args: []string{
		"--web.listen-address=:${NOMAD_PORT_http}",
		"--path.procfs=/host/proc",
		"--path.sysfs=/host/sys",
		"--collector.filesystem.ignored-mount-points",
		"^/(sys|proc|dev|host|etc|rootfs/var/lib/docker/containers|rootfs/var/lib/docker/overlay2|rootfs/run/docker/netns|rootfs/var/lib/docker/aufs)($$|/)",
	},
	Volumes: map[string]string{"/proc": "/host/proc", "/sys": "/host/sys", "/": "/rootfs"},
	Size:    map[string]int{"cpu": 128, "mem": 128},
}

var JenkinsJob = &slomad.Job{
	Name:       "jenkins",
	Image:      getDockerImageString("jenkins"),
	CommonArgs: getCommonJobArgs("docker", "^jenkins-server-0$", 1, 80),
	Ports:      []slomad.Port{{Label: "http", To: 8080}, {Label: "misc", To: 50000}},
	Size:       map[string]int{"cpu": 512, "mem": 1024},
	Storage:    slomad.StringPtr("jenkins"),
	Mounts:     map[string]string{"jenkins-vol": "/var/jenkins_home"},
}

//go:embed config/promtail.yml
var promtailConfig string

var PromtailJob = &slomad.Job{
	Name:       "promtail",
	Image:      getDockerImageString("promtail"),
	JobType:    "system",
	CommonArgs: getCommonJobArgs("docker", "^.*$", 1, 50),
	Env:        map[string]string{"HOSTNAME": "${attr.unique.hostname}"},
	Ports:      []slomad.Port{{Label: "http", To: 3200}},
	Size:       map[string]int{"cpu": 128, "mem": 128},
	Volumes:    map[string]string{"/opt/nomad/data/": "/nomad/", "/data/promtail": "/data"},
	Templates:  map[string]string{"promtail.yml": promtailConfig},
	Args: []string{
		"-config.file=/local/promtail.yml",
		"-server.http-listen-port=${NOMAD_PORT_http}",
	},
}
