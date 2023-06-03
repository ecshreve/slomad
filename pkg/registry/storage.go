package registry

import (
	"fmt"
	"os"

	"github.com/ecshreve/slomad/pkg/slomad"
	nomadApi "github.com/hashicorp/nomad/api"
	"github.com/samsarahq/go/oops"
	log "github.com/sirupsen/logrus"
)

// getStorageArgs returns the common args for the storage controller and node.
//
// TODO: input validation
func getStorageArgs(storage string) []string {
	common := []string{
		"--node-id=${attr.unique.hostname}",
		fmt.Sprintf("--nfs-server=%s", os.Getenv("SYNOLOGY_VAULT")),
		"--mount-options=defaults",
		fmt.Sprintf("--type=%s", storage),
	}

	return common
}

var ControllerJob = slomad.NewStorageJob(slomad.JobParams{
	Name:   "storage-controller",
	Type:   slomad.SERVICE,
	Target: slomad.WORKER,
	TaskConfigParams: slomad.TaskConfigParams{
		Ports: slomad.BasicPorts(0),
		Shape: slomad.DEFAULT_TASK,
		Args:  getStorageArgs("controller"),
	},
})

var NodeJob = slomad.NewStorageJob(slomad.JobParams{
	Name: "storage-node",
	Type: slomad.SYSTEM,
	TaskConfigParams: slomad.TaskConfigParams{
		Ports: slomad.BasicPorts(0),
		Shape: slomad.TINY_TASK,
		Args:  getStorageArgs("node"),
	},
})

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
