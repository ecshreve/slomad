package slomad

import (
	"testing"

	nomadStructs "github.com/hashicorp/nomad/nomad/structs"
	"github.com/stretchr/testify/assert"
)

func TestGetCSIPluginConfig(t *testing.T) {
	testcases := []struct {
		desc     string
		jobName  string
		jobType  JobType
		expected *nomadStructs.TaskCSIPluginConfig
	}{
		{
			desc:     "job type service",
			jobName:  "some-service",
			jobType:  SERVICE,
			expected: nil,
		},
		{
			desc:     "job type system",
			jobName:  "some-system",
			jobType:  SYSTEM,
			expected: nil,
		},
		{
			desc:     "job type batch",
			jobName:  "some-batch",
			jobType:  BATCH,
			expected: nil,
		},
		{
			desc:    "job type storage controller",
			jobName: "storage-controller",
			jobType: STORAGE_CONTROLLER,
			expected: &nomadStructs.TaskCSIPluginConfig{
				ID:       "nfs",
				MountDir: "/csi",
				Type:     nomadStructs.CSIPluginType("controller"),
			},
		},
		{
			desc:    "job type storage node",
			jobName: "storage-node",
			jobType: STORAGE_NODE,
			expected: &nomadStructs.TaskCSIPluginConfig{
				ID:       "nfs",
				MountDir: "/csi",
				Type:     nomadStructs.CSIPluginType("node"),
			},
		},
	}

	for _, tc := range testcases {
		t.Run(tc.desc, func(t *testing.T) {
			j := &Job{Name: tc.jobName, Type: tc.jobType}
			actual := getCSIPluginConfig(j)
			assert.Equal(t, tc.expected, actual)
		})
	}
}
