package registry

import (
	"fmt"
	"os"

	smd "github.com/ecshreve/slomad/pkg/slomad"
)

// ControllerJob is the job for the storage controller.
var ControllerJob = smd.Job{
	Name:   "storage-controller",
	Type:   smd.STORAGE_CONTROLLER,
	Target: smd.NODE0,
	Ports:  smd.BasicPortConfig(0),
	Shape:  smd.DEFAULT_TASK,
	Args:   getStorageArgs("monolith"),
	Priv:   true,
	User:   "root",
}

var NodeJob = smd.Job{
	Name:   "storage-node",
	Type:   smd.STORAGE_NODE,
	Target: smd.NODE,
	Ports:  smd.BasicPortConfig(0),
	Shape:  smd.TINY_TASK,
	Args:   getStorageArgs("node"),
	Priv:   true,
}

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
