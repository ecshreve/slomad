package registry

import "github.com/ecshreve/slomad/pkg/slomad"

var influxPort = 8086
var InfluxJob = &slomad.Job{
	Name:       "influxdb",
	Image:      getDockerImageString("influxdb"),
	Volumes:    map[string]string{"influx_data": "/var/lib/influxdb"},
	CommonArgs: getCommonJobArgs("docker", "^worker-[0-9]+$", 1, 50),
	// Ports:      []slomad.Port{"http", influxPort, tru},
	Size: map[string]int{"cpu": 512, "mem": 512},
}
