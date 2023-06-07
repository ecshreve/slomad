package slomad

import (
	"testing"

	nomadStructs "github.com/hashicorp/nomad/nomad/structs"
	"github.com/stretchr/testify/assert"
)

func TestGetMounts(t *testing.T) {
	testcases := []struct {
		desc     string
		vols     []Volume
		expected []*nomadStructs.VolumeMount
	}{
		{
			desc:     "empty",
			vols:     []Volume{},
			expected: []*nomadStructs.VolumeMount{},
		},
		{
			desc: "no mounts",
			vols: []Volume{
				{Src: "foo", Dst: "bar", Mount: false},
				{Src: "baz", Dst: "qux", Mount: false},
			},
			expected: []*nomadStructs.VolumeMount{},
		},
		{
			desc: "one mount",
			vols: []Volume{
				{Src: "foo", Dst: "bar", Mount: false},
				{Src: "baz", Dst: "qux", Mount: true},
			},
			expected: []*nomadStructs.VolumeMount{
				{Volume: "baz", Destination: "qux", ReadOnly: false},
			},
		},
		{
			desc: "multiple mounts",
			vols: []Volume{
				{Src: "foo", Dst: "bar", Mount: true},
				{Src: "baz", Dst: "qux", Mount: false},
				{Src: "quux", Dst: "corge", Mount: true},
			},
			expected: []*nomadStructs.VolumeMount{
				{Volume: "foo", Destination: "bar", ReadOnly: false},
				{Volume: "quux", Destination: "corge", ReadOnly: false},
			},
		},
	}

	for _, tc := range testcases {
		t.Run(tc.desc, func(t *testing.T) {
			actual := getMounts(tc.vols)
			assert.Equal(t, tc.expected, actual)
		})
	}
}

func TestGetVolumeStrings(t *testing.T) {
	testcases := []struct {
		desc     string
		vols     []Volume
		expected []string
	}{
		{
			desc:     "empty",
			vols:     []Volume{},
			expected: []string{},
		},
		{
			desc: "no mounts",
			vols: []Volume{
				{Src: "foo", Dst: "bar", Mount: false},
				{Src: "baz", Dst: "qux", Mount: false},
			},
			expected: []string{"baz:qux", "foo:bar"},
		},
		{
			desc: "one mount",
			vols: []Volume{
				{Src: "foo", Dst: "bar", Mount: false},
				{Src: "baz", Dst: "qux", Mount: true},
			},
			expected: []string{"foo:bar"},
		},
		{
			desc: "multiple mounts",
			vols: []Volume{
				{Src: "foo", Dst: "bar", Mount: true},
				{Src: "baz", Dst: "qux", Mount: false},
				{Src: "quux", Dst: "corge", Mount: true},
			},
			expected: []string{"baz:qux"},
		},
	}

	for _, tc := range testcases {
		t.Run(tc.desc, func(t *testing.T) {
			actual := getVolumeStrings(tc.vols)
			assert.Equal(t, tc.expected, actual)
		})
	}
}

func TestGetNomadVolumeReq(t *testing.T) {
	testcases := []struct {
		desc     string
		vols     []Volume
		expected map[string]*nomadStructs.VolumeRequest
	}{
		{
			desc:     "empty",
			vols:     []Volume{},
			expected: map[string]*nomadStructs.VolumeRequest{},
		},
		{
			desc: "no mounts",
			vols: []Volume{
				{Src: "foo", Dst: "bar", Mount: false},
				{Src: "baz", Dst: "qux", Mount: false},
			},
			expected: map[string]*nomadStructs.VolumeRequest{},
		},
		{
			desc: "one mount",
			vols: []Volume{
				{Src: "foo", Dst: "bar", Mount: false},
				{Src: "baz", Dst: "qux", Mount: true},
			},
			expected: map[string]*nomadStructs.VolumeRequest{
				"baz": {
					Name:           "baz",
					Source:         "baz",
					Type:           "csi",
					ReadOnly:       false,
					AccessMode:     "single-node-writer",
					AttachmentMode: "file-system",
				},
			},
		},
	}

	for _, tc := range testcases {
		t.Run(tc.desc, func(t *testing.T) {
			actual := getNomadVolumeReq(tc.vols)
			assert.Equal(t, tc.expected, actual)
		})
	}
}

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
