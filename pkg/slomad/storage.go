package slomad

import (
	"fmt"
	"sort"

	nomadStructs "github.com/hashicorp/nomad/nomad/structs"
)

type Volume struct {
	Src   string
	Dst   string
	Mount bool
}

func NewDockerVolume(src, dst string) *Volume {
	return &Volume{
		Src: src,
		Dst: dst,
	}
}

func NewNomadVolume(src, dst string) *Volume {
	return &Volume{
		Src:   src,
		Dst:   dst,
		Mount: true,
	}
}

// getMounts converts a list of Volumes to a list of Nomad VolumeMounts.
func getMounts(vols []Volume) []*nomadStructs.VolumeMount {
	volMounts := []*nomadStructs.VolumeMount{}
	for _, vol := range vols {
		if !vol.Mount {
			continue
		}

		volMounts = append(volMounts, &nomadStructs.VolumeMount{
			Volume:      vol.Src,
			Destination: vol.Dst,
			ReadOnly:    false,
		})
	}
	return volMounts
}

// getNomadVolumes converts a list of Volumes to a list of Nomad Volumes.
func getNomadVolumes(storage string) map[string]*nomadStructs.VolumeRequest {
	if storage == "" || storage == "controller" || storage == "node" {
		return nil
	}

	csiVols := map[string]*nomadStructs.VolumeRequest{}
	volName := fmt.Sprintf("%s-vol", storage)
	csiVols[volName] = &nomadStructs.VolumeRequest{
		Name:   volName,
		Source: volName,

		Type:           "csi",
		ReadOnly:       false,
		AccessMode:     "single-node-writer",
		AttachmentMode: "file-system",
	}

	return csiVols
}

// toVolumeString converts a list of Volumes to a list of Volume strings.
// These are meant to be passed to the docker driver.
func toVolumeStrings(vols []Volume) []string {
	volStrings := []string{}
	for _, vol := range vols {
		if vol.Mount {
			continue
		}

		volStrings = append(volStrings, fmt.Sprintf("%s:%s", vol.Src, vol.Dst))
	}

	sort.Strings(volStrings)
	return volStrings
}
