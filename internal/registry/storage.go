package registry

import (
	"fmt"
	"os"

	smd "github.com/ecshreve/slomad/pkg/slomad"
)

var ControllerJob = smd.NewAppJob(smd.JobParams{
	Name:   "storage-controller",
	Type:   smd.STORAGE_CONTROLLER,
	Target: smd.WORKER,
	TaskConfigParams: smd.TaskConfigParams{
		Ports: smd.BasicPortConfig(0),
		Shape: smd.DEFAULT_TASK,
		Args:  getStorageArgs("controller"),
	},
})

var NodeJob = smd.NewAppJob(smd.JobParams{
	Name: "storage-node",
	Type: smd.STORAGE_NODE,
	TaskConfigParams: smd.TaskConfigParams{
		Ports: smd.BasicPortConfig(0),
		Shape: smd.TINY_TASK,
		Args:  getStorageArgs("node"),
	},
})

// getStorageArgs returns the common args for the storage controller and node.
//
// TODO: input validation
func getStorageArgs(storage string) []string {
	common := []string{
		"--node-id=${attr.unique.hostname}",
		fmt.Sprintf("--nfs-server=%s", os.Getenv("NFS_MOUNT")),
		"--mount-options=defaults",
		fmt.Sprintf("--type=%s", storage),
	}

	return common
}
