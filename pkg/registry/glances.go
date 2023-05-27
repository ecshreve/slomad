package registry

import (
	"fmt"
	"time"

	nomadStructs "github.com/hashicorp/nomad/nomad/structs"
	"github.com/samsarahq/go/oops"
	log "github.com/sirupsen/logrus"
)

// var GlancesJob = &slomad.Job{
// 	Name:       "glances",
// 	Image:      getDockerImageString("glances"),
// 	JobType:    "system",
// 	CommonArgs: getCommonJobArgs("docker", "^coderbox*$", 1, 50),
// 	Env:        map[string]string{"GLANCES_OPT": "-s"},
// 	Ports:      []slomad.Port{{Label: "http", To: 61209, From: slomad.IntPtr(61209)}},
// 	Size:       map[string]int{"cpu": 128, "mem": 128},
// 	Volumes:    map[string]string{"/var/run/docker.sock": "/var/run/docker.sock"},
// }

var GlancesJob = nomadStructs.Job{
	ID:          "glances",
	Name:        "glances",
	Region:      "global",
	Priority:    80,
	Datacenters: []string{"dcs"},
	Type:        "system",
	TaskGroups: []*nomadStructs.TaskGroup{
		{
			Name:  "glances",
			Count: 1,
			Tasks: []*nomadStructs.Task{
				{
					Name:   "glances",
					Driver: "exec",
					Config: map[string]interface{}{
						"command": "glances",
						"args":    []string{"-C", "/etc/glances/glances.conf", "--quiet", "--export", "influxdb2"},
					},
					Resources: &nomadStructs.Resources{
						CPU:      128,
						MemoryMB: 128,
					},
					LogConfig: nomadStructs.DefaultLogConfig(),
					Services: []*nomadStructs.Service{
						{
							Name:      "glances",
							PortLabel: "http",
							TaskName:  "glances",
							Checks: []*nomadStructs.ServiceCheck{
								{
									Name:          fmt.Sprintf("%s = tcp check", "glances"),
									Type:          nomadStructs.ServiceCheckTCP,
									Interval:      10 * time.Second,
									Timeout:       2 * time.Second,
									InitialStatus: "passing",
								},
							},
							Provider: "consul",
						},
					},
				},
			},
			RestartPolicy: &nomadStructs.DefaultServiceJobRestartPolicy,
			EphemeralDisk: &nomadStructs.EphemeralDisk{
				SizeMB: 128,
			},
			Networks: []*nomadStructs.NetworkResource{
				{
					Mode: "host",
					ReservedPorts: []nomadStructs.Port{
						{
							Label: "http",
							Value: 61209,
							To:    61209,
						},
					},
				},
			},
		},
	},
	Namespace: "default",
	Constraints: []*nomadStructs.Constraint{
		{
			LTarget: "${attr.unique.hostname}",
			RTarget: "^.*$",
			Operand: "regexp",
		},
	},
}

func DeployGlances() error {
	job := &GlancesJob
	if err := job.Validate(); err != nil {
		log.Errorf("Nomad job validation failed. Error: %s\n", err)
		return err
	}

	apiJob, err := convertJob(job)
	if err != nil {
		log.Errorf("Failed to convert nomad job in api call. Error: %s\n", err)
		return err
	}

	if err = submitApiJob(apiJob); err != nil {
		return oops.Wrapf(err, "error submitting api job")
	}

	return nil
}
