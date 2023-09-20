package registry

import (
	smd "github.com/ecshreve/slomad/pkg/slomad"
)

// PlexJob is the job definition for the plex service.
//
// @DEPRECATED: This job is no longer used.
var PlexJob = smd.Job{
	Name:   "plex",
	Type:   smd.SERVICE,
	Target: smd.PLEXBOX,
	Ports:  []*smd.Port{{Label: "http", To: 32400, From: 32400, Static: true}},
	Shape:  smd.PLEX_TASK,
	User:   "root",
	Env: map[string]string{
		"TZ":           "America/Los_Angeles",
		"VERSION":      "docker",
		"ADVERTISE_IP": "http://plex.slab.lan:80",
		"PGID":         "100",
		"PUID":         "1027",
	},
	Volumes: []smd.Volume{
		{Src: "/mnt/nfs/config/plex", Dst: "/config"},
		{Src: "/mnt/nfs/media/music", Dst: "/music"},
		{Src: "/mnt/nfs/media/tv", Dst: "/tv"},
		{Src: "/mnt/nfs/media/movies", Dst: "/movies"},
		{Src: "/dev/shm", Dst: "/transcode"},
	},
}
