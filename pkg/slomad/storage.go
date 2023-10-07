package slomad

import (
	"fmt"
	"sort"

	nomadApi "github.com/hashicorp/nomad/api"
	nomadStructs "github.com/hashicorp/nomad/nomad/structs"
	"github.com/samsarahq/go/oops"
	log "github.com/sirupsen/logrus"
)

type Volume struct {
	Src   string
	Dst   string
	Mount bool
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

// getVolumeString converts a list of Volumes to a list of Volume strings.
// These are meant to be passed to the docker driver.
func getVolumeStrings(vols []Volume) []string {
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

// getNomadVolumeReq converts a slice of Volumes to map of nomad VolumeRequest.
//
// TODO: validation
func getNomadVolumeReq(vols []Volume) map[string]*nomadStructs.VolumeRequest {
	csiVols := map[string]*nomadStructs.VolumeRequest{}
	for _, v := range vols {
		if !v.Mount {
			continue
		}

		volName := v.Src
		csiVols[volName] = &nomadStructs.VolumeRequest{
			Name:   volName,
			Source: volName,

			Type:           "csi",
			ReadOnly:       false,
			AccessMode:     "single-node-writer",
			AttachmentMode: "file-system",
		}
	}

	return csiVols
}

// getCSIPluginConfig returns a CSIPluginConfig for a given job.
func getCSIPluginConfig(j *Job) *nomadStructs.TaskCSIPluginConfig {
	if j.Type != STORAGE_CONTROLLER && j.Type != STORAGE_NODE && j.Type != STORAGE_MONOLITH {
		log.Info("job type is not storage, skipping CSIPluginConfig")
		return nil
	}

	// storageType := strings.Split(j.Name, "-")[1]
	return &nomadStructs.TaskCSIPluginConfig{
		ID:       "nfs",
		MountDir: "/csi",
		Type:     nomadStructs.CSIPluginTypeMonolith,
	}
}

func CreateVolume(volName string) error {
	nomadConfig := nomadApi.DefaultConfig()
	nomadClient, err := nomadApi.NewClient(nomadConfig)
	if err != nil {
		return oops.Wrapf(err, "unable to create nomad api client")
	}

	vol := &nomadApi.CSIVolume{
		Name:     fmt.Sprintf("%s-vol", volName),
		ID:       fmt.Sprintf("%s-vol", volName),
		PluginID: "nfs",
		Capacity: 100000,
		RequestedCapabilities: []*nomadApi.CSIVolumeCapability{
			{
				AccessMode:     nomadApi.CSIVolumeAccessMode("single-node-writer"),
				AttachmentMode: nomadApi.CSIVolumeAttachmentMode("file-system"),
			},
		},
	}

	nomadVol, _, nomadErr := nomadClient.CSIVolumes().Create(vol, nil)
	if nomadErr != nil {
		return oops.Wrapf(nomadErr, "error creating volume: %+v", vol)
	}

	log.Infof("Sucessfully created nomad volume: %v\n", nomadVol)
	return nil
}
