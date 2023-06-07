package slomad

import (
	"testing"

	nomadStructs "github.com/hashicorp/nomad/nomad/structs"
	"github.com/samsarahq/go/snapshotter"
	"github.com/stretchr/testify/assert"
)

func TestGetService(t *testing.T) {
	snap := snapshotter.New(t)
	defer snap.Verify()

	storageSrv := getService("storage-controller", "http")
	snap.Snapshot("storage-controller", storageSrv)

	basicSrv := getService("basic", "http")
	snap.Snapshot("basic", basicSrv)
}

func TestGetTemplates(t *testing.T) {
	testcases := []struct {
		desc     string
		inp      map[string]string
		expected []*nomadStructs.Template
	}{
		{
			desc:     "empty template",
			inp:      map[string]string{},
			expected: nil,
		},
		{
			desc: "basic template",
			inp: map[string]string{
				"test": "test",
			},
			expected: []*nomadStructs.Template{
				{
					EmbeddedTmpl: "test",
					DestPath:     "local/config/test",
					ChangeMode:   "signal",
					ChangeSignal: "SIGHUP",
				},
			},
		},
	}

	for _, tc := range testcases {
		t.Run(tc.desc, func(t *testing.T) {
			assert.Equal(t, tc.expected, getTemplates(tc.inp))
		})
	}
}

func TestGetConfig(t *testing.T) {
	testcases := []struct {
		desc     string
		img      string
		jt       JobType
		args     []string
		ports    []*Port
		vols     []Volume
		expected map[string]interface{}
	}{
		{
			desc:  "basic config",
			img:   "test-image",
			jt:    SERVICE,
			args:  []string{"echo hello"},
			ports: []*Port{{Label: "http", To: 8080, From: 0, Static: false}},
			vols:  []Volume{{Src: "/tmp", Dst: "/tmp", Mount: false}},
			expected: map[string]interface{}{
				"image":   "test-image",
				"args":    []string{"echo hello"},
				"ports":   []string{"http"},
				"volumes": []string{"/tmp:/tmp"},
			},
		},
		{
			desc:  "storage controller config",
			img:   "storage-image",
			jt:    STORAGE_CONTROLLER,
			args:  []string{"echo hello"},
			ports: []*Port{{Label: "http", To: 8080, From: 0, Static: false}},
			vols:  []Volume{{Src: "/tmp", Dst: "/tmp", Mount: false}},
			expected: map[string]interface{}{
				"image":        "storage-image",
				"args":         []string{"echo hello"},
				"ports":        []string{"http"},
				"volumes":      []string{"/tmp:/tmp"},
				"privileged":   true,
				"network_mode": "host",
			},
		},
		{
			desc:  "storage node config",
			img:   "storage-image",
			jt:    STORAGE_NODE,
			args:  []string{"echo hello"},
			ports: []*Port{{Label: "http", To: 8080, From: 0, Static: false}},
			vols:  []Volume{{Src: "/tmp", Dst: "/tmp", Mount: false}},
			expected: map[string]interface{}{
				"image":        "storage-image",
				"args":         []string{"echo hello"},
				"ports":        []string{"http"},
				"volumes":      []string{"/tmp:/tmp"},
				"privileged":   true,
				"network_mode": "host",
			},
		},
	}

	for _, tc := range testcases {
		t.Run(tc.desc, func(t *testing.T) {
			actual := getConfig(tc.img, tc.jt, tc.args, tc.ports, tc.vols)
			assert.EqualValues(t, tc.expected, actual)
		})
	}
}

func TestGetDisk(t *testing.T) {
	expected := &nomadStructs.EphemeralDisk{
		SizeMB: 500,
	}

	actual := getDisk()
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
