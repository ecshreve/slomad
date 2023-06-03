package registry

import (
	"fmt"
)

func getDockerImageString(name string) *string {
	imageStr := fmt.Sprintf("registry.slab.lan:5000/%s:custom", name)
	return &imageStr
}
