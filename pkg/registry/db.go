package registry

import "github.com/ecshreve/slomad/pkg/slomad"

var postgresPort = 5432
var PostgresJob = &slomad.Job{
	Name:  "coderdb",
	Image: getDockerImageString("postgres"),
	Env: map[string]string{
		"POSTGRES_USER":     "coder-user",
		"POSTGRES_PASSWORD": "password",
		"POSTGRES_DB":       "basecoder",
	},
	Volumes:    map[string]string{"coder_data": "/var/lib/postgresql/data"},
	CommonArgs: getCommonJobArgs("docker", "^nuck$", 1, 50),
	Ports:      []slomad.Port{*slomad.NewPort("http", postgresPort, &postgresPort)},
	Size:       map[string]int{"cpu": 1024, "mem": 1024},
}

var AdminerJob = &slomad.Job{
	Name:  "adminer",
	Image: getDockerImageString("adminer"),
	Env: map[string]string{
		"ADMINER_DEFAULT_SERVER": "postgres://coder-user:password@coderdb:5432/basecoder",
	},
	CommonArgs: getCommonJobArgs("docker", "^nuck$", 1, 50),
	Ports:      []slomad.Port{*slomad.NewPort("http", 8080, nil)},
	Size:       map[string]int{"cpu": 256, "mem": 256},
}

var influxPort = 8086
var InfluxJob = &slomad.Job{
	Name:       "influxdb",
	Image:      getDockerImageString("influxdb"),
	Volumes:    map[string]string{"influx_data": "/var/lib/influxdb"},
	CommonArgs: getCommonJobArgs("docker", "^worker-[0-9]+$", 1, 50),
	Ports:      []slomad.Port{*slomad.NewPort("http", influxPort, &influxPort)},
	Size:       map[string]int{"cpu": 512, "mem": 512},
}
