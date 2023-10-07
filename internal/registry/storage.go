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
	Args:   buildStorageArgs("monolith"),
	Priv:   true,
	User:   "root",
}

// buildStorageArgs returns the args for the storage related jobs.
//
// Currently this is running in monolith mode, so the storage controller
// and node are running on the same host.
func buildStorageArgs(storage string) []string {
	return []string{
		"--node-id=${attr.unique.hostname}",
		fmt.Sprintf("--nfs-server=%s", os.Getenv("NFS_MOUNT")),
		"--mount-options=defaults",
		fmt.Sprintf("--type=%s", storage),
	}
}

// DEPRECATED: This job is no longer used.
// var NodeJob = smd.Job{
// 	Name:   "storage-node",
// 	Type:   smd.STORAGE_NODE,
// 	Target: smd.NODE,
// 	Ports:  smd.BasicPortConfig(0),
// 	Shape:  smd.TINY_TASK,
// 	Args:   getStorageArgs("node"),
// 	Priv:   true,
// }
