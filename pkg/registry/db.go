package registry

import "github.com/ecshreve/slomad/pkg/slomad"

// TODO: mount nomad volume and persist data
var InfluxDBJob = slomad.NewAppJob(slomad.JobParams{
	Name:   "influxdb",
	Type:   slomad.SERVICE,
	Target: slomad.WORKER,
	TaskConfigParams: slomad.TaskConfigParams{
		Ports: slomad.BasicPortConfig(8086),
		Shape: slomad.LARGE_TASK,
	},
	StorageParams: slomad.StorageParams{
		Volumes: []slomad.Volume{
			{Src: "influx_data", Dst: "/var/lib/influxdb"},
		},
	},
})
