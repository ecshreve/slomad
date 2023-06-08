package slomad

import (
	"testing"

	"github.com/samsarahq/go/snapshotter"
)

func TestGetService(t *testing.T) {
	snap := snapshotter.New(t)
	defer snap.Verify()

	storageSrv := getService("storage-controller", "http")
	snap.Snapshot("storage-controller", storageSrv)

	basicSrv := getService("basic", "http")
	snap.Snapshot("basic", basicSrv)
}
