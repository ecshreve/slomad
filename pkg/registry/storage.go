package registry

import (
	"fmt"
	"os"

	"github.com/ecshreve/slomad/pkg/slomad"
	nomadApi "github.com/hashicorp/nomad/api"
	"github.com/samsarahq/go/oops"
	log "github.com/sirupsen/logrus"
)

func getStorageArgs(storage string) []string {
	common := []string{
		"--node-id=${attr.unique.hostname}",
		fmt.Sprintf("--nfs-server=%s", os.Getenv("SYNOLOGY_VAULT")),
		"--mount-options=defaults",
	}

	switch storage {
	case "controller":
		return append(common, "--type=controller")
	case "node":
		return append(common, "--type=node")
	default:
		log.Fatalf("unknown storage type: %s", storage)
		return nil
	}
}

var ControllerJobOld = &slomad.Job{
	Name:       "storage-controller",
	Image:      getDockerImageString("csi-nfs-plugin"),
	CommonArgs: getCommonJobArgs("docker", "^worker-0$", 1, 99),
	Storage:    slomad.StringPtr("controller"),
	Ports:      []slomad.Port{{Label: "http"}},
	Args:       getStorageArgs("controller"),
	Size:       map[string]int{"cpu": 512, "mem": 512},
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

var NodeJob = &slomad.Job{
	Name:       "storage-node",
	Image:      getDockerImageString("csi-nfs-plugin"),
	CommonArgs: getCommonJobArgs("docker", "^.*$", 1, 98),
	JobType:    "system",
	Storage:    slomad.StringPtr("node"),
	Ports:      []slomad.Port{{Label: "http"}},
	Args:       getStorageArgs("node"),
	Size:       map[string]int{"cpu": 128, "mem": 128},
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
