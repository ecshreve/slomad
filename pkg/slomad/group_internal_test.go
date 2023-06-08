package slomad

import (
	"testing"

	nomadStructs "github.com/hashicorp/nomad/nomad/structs"
	"github.com/stretchr/testify/assert"
)

func TestGetDisk(t *testing.T) {
	expected := &nomadStructs.EphemeralDisk{
		SizeMB: 500,
	}

	actual := getDisk(500)
	assert.Equal(t, expected, actual)
}

func TestGetConstraint(t *testing.T) {
	testcases := []struct {
		desc     string
		dt       DeployTarget
		expected *nomadStructs.Constraint
	}{
		{
			desc: "deploy target all",
			dt:   ALL,
			expected: &nomadStructs.Constraint{
				LTarget: "${attr.unique.hostname}",
				RTarget: "^.*$",
				Operand: "regexp",
			},
		},
		{
			desc: "deploy target workers",
			dt:   WORKER,
			expected: &nomadStructs.Constraint{
				LTarget: "${attr.unique.hostname}",
				RTarget: "^worker-[0-9]+$",
				Operand: "regexp",
			},
		},
		{
			desc: "deploy target devbox",
			dt:   DEVBOX,
			expected: &nomadStructs.Constraint{
				LTarget: "${attr.unique.hostname}",
				RTarget: "^devbox$",
				Operand: "regexp",
			},
		},
	}

	for _, tc := range testcases {
		t.Run(tc.desc, func(t *testing.T) {
			actual := getConstraint(tc.dt)
			assert.Equal(t, tc.expected, actual)
		})
	}
}

func TestGetReschedulePolicy(t *testing.T) {
	testcases := []struct {
		desc     string
		jobType  JobType
		expected *nomadStructs.ReschedulePolicy
	}{
		{
			desc:     "job type service",
			jobType:  SERVICE,
			expected: &nomadStructs.DefaultServiceJobReschedulePolicy,
		},
		{
			desc:     "job type system",
			jobType:  SYSTEM,
			expected: nil,
		},
		{
			desc:     "job type batch",
			jobType:  BATCH,
			expected: nil,
		},
		{
			desc:     "job type storage controller",
			jobType:  STORAGE_CONTROLLER,
			expected: &nomadStructs.DefaultServiceJobReschedulePolicy,
		},
		{
			desc:     "job type storage node",
			jobType:  STORAGE_NODE,
			expected: nil,
		},
	}

	for _, tc := range testcases {
		t.Run(tc.desc, func(t *testing.T) {
			actual := getReschedulePolicy(tc.jobType)
			assert.Equal(t, tc.expected, actual)
		})
	}
}
