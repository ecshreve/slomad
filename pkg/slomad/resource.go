package slomad

import (
	nomadStructs "github.com/hashicorp/nomad/nomad/structs"
)

type ResourceValue int

const (
	TINY     ResourceValue = 128
	SMALL    ResourceValue = 256
	DEFAULT  ResourceValue = 512
	LARGE    ResourceValue = 1024
	XLARGE   ResourceValue = 2048
	XXLARGE  ResourceValue = 4096
	XXXLARGE ResourceValue = 8192
)

// TaskResource is a struct that represents the CPU and MEM resources for a task.
type TaskResource struct {
	CPU ResourceValue
	MEM ResourceValue
}

var (
	DEFAULT_TASK = TaskResource{CPU: DEFAULT, MEM: DEFAULT}
	TINY_TASK    = TaskResource{CPU: TINY, MEM: TINY}
	SMALL_TASK   = TaskResource{CPU: SMALL, MEM: SMALL}
	LARGE_TASK   = TaskResource{CPU: LARGE, MEM: LARGE}
	XLARGE_TASK  = TaskResource{CPU: XLARGE, MEM: XLARGE}
	MEM_TASK     = TaskResource{CPU: SMALL, MEM: LARGE}
	COMPUTE_TASK = TaskResource{CPU: LARGE, MEM: SMALL}
	PLEX_TASK    = TaskResource{CPU: LARGE, MEM: XXLARGE}
)

// TODO: remove this function?
//
// CustomTaskResource is a helper function to create a TaskResource with custom values.
// func CustomTaskResource(cpu ResourceValue, mem ResourceValue) TaskResource {
// 	return TaskResource{CPU: cpu, MEM: mem}
// }

// getResource is a helper function to convert a TaskResource to a Nomad Resources struct.
func getResource(tr TaskResource) *nomadStructs.Resources {
	return &nomadStructs.Resources{
		CPU:      int(tr.CPU),
		MemoryMB: int(tr.MEM),
	}
}
