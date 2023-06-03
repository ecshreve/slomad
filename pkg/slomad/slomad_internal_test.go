package slomad

import (
	"testing"

	nomadStructs "github.com/hashicorp/nomad/nomad/structs"
	"github.com/stretchr/testify/assert"
)

func TestGetResource(t *testing.T) {
	testcases := []struct {
		desc     string
		tr       TaskResource
		expected *nomadStructs.Resources
	}{
		{
			desc:     "default task",
			tr:       DEFAULT_TASK,
			expected: &nomadStructs.Resources{CPU: 512, MemoryMB: 512},
		},
		{
			desc:     "tiny task",
			tr:       TINY_TASK,
			expected: &nomadStructs.Resources{CPU: 128, MemoryMB: 128},
		},
		{
			desc:     "small task",
			tr:       SMALL_TASK,
			expected: &nomadStructs.Resources{CPU: 256, MemoryMB: 256},
		},
		{
			desc:     "large task",
			tr:       LARGE_TASK,
			expected: &nomadStructs.Resources{CPU: 1024, MemoryMB: 1024},
		},
		{
			desc:     "xlarge task",
			tr:       XLARGE_TASK,
			expected: &nomadStructs.Resources{CPU: 2048, MemoryMB: 2048},
		},
		{
			desc:     "mem task",
			tr:       MEM_TASK,
			expected: &nomadStructs.Resources{CPU: 256, MemoryMB: 1024},
		},
		{
			desc:     "compute task",
			tr:       COMPUTE_TASK,
			expected: &nomadStructs.Resources{CPU: 1024, MemoryMB: 256},
		},
	}

	for _, tc := range testcases {
		t.Run(tc.desc, func(t *testing.T) {
			assert.Equal(t, tc.expected, getResource(tc.tr))
		})
	}
}
