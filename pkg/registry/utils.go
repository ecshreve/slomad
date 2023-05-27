package registry

import (
	"fmt"

	"github.com/ecshreve/slomad/pkg/slomad"
)

func getDockerImageString(name string) *string {
	imageStr := fmt.Sprintf("registry.slab.lan:5000/%s:custom", name)
	return &imageStr
}

func getCommonJobArgs(driver, constraint string, c, p int) slomad.CommonArgs {
	return slomad.CommonArgs{
		Driver:     driver,
		Constraint: constraint,
		Count:      c,
		Priority:   p,
	}
}
