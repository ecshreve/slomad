package registry

import (
	"fmt"
	"os"
)

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
