<!-- Code generated by gomarkdoc. DO NOT EDIT -->

# registry

```go
import "github.com/ecshreve/slomad/internal/registry"
```

## Index

- [Variables](<#variables>)
- [func DeployTraefikJob(confirm bool) error](<#func-deploytraefikjob>)
- [func getStorageArgs(storage string) []string](<#func-getstorageargs>)
- [func promConfigHelper(tmpl string) string](<#func-promconfighelper>)


## Variables

```go
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
```

```go
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
        Storage: utils.StringPtr("grafana"),
        Volumes: []smd.Volume{{Src: "grafana-vol", Dst: "/var/lib/grafana", Mount: true}},
    },
})
```

TODO: mount nomad volume and persist data

```go
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
```

```go
var LokiJob = smd.NewAppJob(smd.JobParams{
    Name:   "loki",
    Type:   smd.SERVICE,
    Target: smd.WORKER,
    TaskConfigParams: smd.TaskConfigParams{
        Ports: smd.BasicPortConfig(3100),
        Shape: smd.TINY_TASK,
    },
})
```

```go
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
```

```go
var NodeJob = smd.NewStorageJob(smd.JobParams{
    Name: "storage-node",
    Type: smd.SYSTEM,
    TaskConfigParams: smd.TaskConfigParams{
        Ports: smd.BasicPortConfig(0),
        Shape: smd.TINY_TASK,
        Args:  getStorageArgs("node"),
    },
})
```

PrometheusJob is a Job for the Prometheus service.

```go
var PrometheusJob = smd.NewAppJob(smd.JobParams{
    Name:   "prometheus",
    Type:   smd.SERVICE,
    Target: smd.WORKER,
    TaskConfigParams: smd.TaskConfigParams{
        Ports:     smd.BasicPortConfig(9090),
        Shape:     smd.LARGE_TASK,
        Templates: map[string]string{"prometheus.yml": promConfigHelper(prometheusConfig)},
    },
    StorageParams: smd.StorageParams{
        Storage: utils.StringPtr("prometheus"),
        Volumes: []smd.Volume{{Src: "local/config", Dst: "/etc/prometheus"}},
    },
})
```

```go
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
        Templates: map[string]string{"promtail.yml": promConfigHelper(promtailConfig)},
    },
    StorageParams: smd.StorageParams{
        Volumes: []smd.Volume{
            {Src: "/opt/nomad/data/", Dst: "/nomad/"},
            {Src: "/data/promtail", Dst: "/data"},
        },
    },
})
```

```go
var SpeedtestJob = smd.NewAppJob(smd.JobParams{
    Name:   "speedtest",
    Type:   smd.SERVICE,
    Target: smd.WORKER,
    TaskConfigParams: smd.TaskConfigParams{
        Ports: smd.BasicPortConfig(80),
        Shape: smd.TINY_TASK,
    },
})
```

```go
var TraefikJob = nomadStructs.Job{
    ID:          "traefik",
    Name:        "traefik",
    Region:      "global",
    Priority:    92,
    Datacenters: []string{"dcs"},
    Type:        "service",
    TaskGroups: []*nomadStructs.TaskGroup{
        {
            Services: []*nomadStructs.Service{
                {
                    Name:      "traefik-web",
                    PortLabel: "web",
                    Checks: []*nomadStructs.ServiceCheck{
                        {
                            Name:          fmt.Sprintf("%s = tcp check", "traefik-web"),
                            Type:          nomadStructs.ServiceCheckTCP,
                            Interval:      10 * time.Second,
                            Timeout:       2 * time.Second,
                            InitialStatus: "passing",
                        },
                    },
                    Provider: "consul",
                },
            },
            Name:  "traefik",
            Count: 1,
            Tasks: []*nomadStructs.Task{
                {
                    Name:   "traefik",
                    Driver: "docker",
                    Config: map[string]interface{}{
                        "image":        "reg.slab.lan:5000/traefik:latest",
                        "network_mode": "host",
                        "args": []string{
                            "--entryPoints.web.address=:80",
                            "--entryPoints.websecure.address=:443",
                            "--entryPoints.admin.address=:8081",
                            "--entrypoints.websecure.http.redirections.entryPoint.to=web",
                            "--entrypoints.websecure.http.redirections.entryPoint.scheme=http",
                            "--accesslog=true",
                            "--api=true",
                            "--api.dashboard=true",
                            "--api.insecure=true",
                            "--ping=true",
                            "--ping.entryPoint=admin",
                            "--providers.consulcatalog=true",
                            "--providers.consulcatalog.endpoint.address=10.35.220.50:8500",
                            "--providers.consulcatalog.prefix=traefik",
                            "--providers.consulcatalog.refreshInterval=30s",
                            "--providers.consulcatalog.exposedByDefault=false",
                            "--providers.consulcatalog.defaultrule=Host(`{{ .Name }}.slabstaging.lan`)",
                            "--providers.consulcatalog.endpoint.tls.insecureskipverify=true",
                        },
                    },
                    Resources: &nomadStructs.Resources{
                        CPU:      512,
                        MemoryMB: 512,
                    },
                    LogConfig: nomadStructs.DefaultLogConfig(),
                    Services: []*nomadStructs.Service{
                        {
                            Name:      "traefik",
                            PortLabel: "websecure",
                            Tags: []string{
                                "traefik.enable=true",
                                "traefik.http.routers.api.rule=Host(`traefik.slabstaging.lan`)",
                                "traefik.http.routers.api.service=api@internal",
                            },
                            TaskName: "traefik",
                            Checks: []*nomadStructs.ServiceCheck{
                                {
                                    Name:          fmt.Sprintf("%s = http check", "traefik"),
                                    Type:          nomadStructs.ServiceCheckHTTP,
                                    Interval:      10 * time.Second,
                                    Timeout:       2 * time.Second,
                                    InitialStatus: "passing",
                                    Path:          "/ping",
                                    PortLabel:     "admin",
                                    TaskName:      "traefik",
                                },
                            },
                            Provider: "consul",
                        },
                    },
                },
            },
            RestartPolicy:    nomadStructs.NewRestartPolicy("service"),
            ReschedulePolicy: &nomadStructs.DefaultServiceJobReschedulePolicy,
            EphemeralDisk: &nomadStructs.EphemeralDisk{
                SizeMB: 256,
            },
            Networks: []*nomadStructs.NetworkResource{
                {
                    Mode: "host",
                    ReservedPorts: []nomadStructs.Port{
                        {
                            Label: "web",
                            Value: 80,
                            To:    0,
                        },
                        {
                            Label: "websecure",
                            Value: 443,
                            To:    0,
                        },
                        {
                            Label: "admin",
                            Value: 8081,
                            To:    0,
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
            RTarget: "worker-0",
            Operand: "regexp",
        },
    },
}
```

```go
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
```

```go
var prometheusConfig string
```

```go
var promtailConfig string
```

## func DeployTraefikJob

```go
func DeployTraefikJob(confirm bool) error
```

DeployTraefikJob deploys the Traefik job to Nomad.

## func getStorageArgs

```go
func getStorageArgs(storage string) []string
```

getStorageArgs returns the common args for the storage controller and node.

TODO: input validation

## func promConfigHelper

```go
func promConfigHelper(tmpl string) string
```



Generated by [gomarkdoc](<https://github.com/princjef/gomarkdoc>)