package registry

import "github.com/ecshreve/slomad/pkg/slomad"

var CoderJob = &slomad.Job{
	Name:  "coder",
	Image: getDockerImageString("coder"),
	Env: map[string]string{
		"CODER_ACCESS_URL":          "http://coder.slab.lan",
		"CODER_HTTP_ADDRESS":        "0.0.0.0:8080",
		"CODER_TLS_ENABLE":          "false",
		"CODER_PG_CONNECTION_URL":   "postgresql://coder-user@coderdb.service.consul:5432/basecoder?sslmode=disable&password=password",
	},
	Caps: 		 []string{"NET_ADMIN"},
	Volumes:    map[string]string{"/var/run/docker.sock": "/var/run/docker.sock"},
	CommonArgs: getCommonJobArgs("docker", "^nuck$", 1, 90),
	Ports:      []slomad.Port{*slomad.NewPort("http", 8080, nil)},
	Size:       map[string]int{"cpu": 1024, "mem": 2048},
}