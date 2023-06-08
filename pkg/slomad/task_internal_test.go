package slomad

import (
	"testing"

	nomadStructs "github.com/hashicorp/nomad/nomad/structs"
	"github.com/stretchr/testify/assert"
)

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
		name     string
		jt       JobType
		args     []string
		ports    []*Port
		vols     []Volume
		expected map[string]interface{}
	}{
		{
			desc:  "basic config",
			name:  "test",
			jt:    SERVICE,
			args:  []string{"echo hello"},
			ports: []*Port{{Label: "http", To: 8080, From: 0, Static: false}},
			vols:  []Volume{{Src: "/tmp", Dst: "/tmp", Mount: false}},
			expected: map[string]interface{}{
				"image":   "reg.slab.lan:5000/test",
				"args":    []string{"echo hello"},
				"ports":   []string{"http"},
				"volumes": []string{"/tmp:/tmp"},
			},
		},
		{
			desc:  "storage controller config",
			name:  "storage-controller",
			jt:    STORAGE_CONTROLLER,
			args:  []string{"echo hello"},
			ports: []*Port{{Label: "http", To: 8080, From: 0, Static: false}},
			vols:  []Volume{{Src: "/tmp", Dst: "/tmp", Mount: false}},
			expected: map[string]interface{}{
				"image":        "reg.slab.lan:5000/csi-nfs-plugin",
				"args":         []string{"echo hello"},
				"ports":        []string{"http"},
				"volumes":      []string{"/tmp:/tmp"},
				"privileged":   true,
				"network_mode": "host",
			},
		},
		{
			desc:  "storage node config",
			name:  "storage-node",
			jt:    STORAGE_NODE,
			args:  []string{"echo hello"},
			ports: []*Port{{Label: "http", To: 8080, From: 0, Static: false}},
			vols:  []Volume{{Src: "/tmp", Dst: "/tmp", Mount: false}},
			expected: map[string]interface{}{
				"image":        "reg.slab.lan:5000/csi-nfs-plugin",
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
			actual := getConfig(tc.name, tc.jt, tc.args, tc.ports, tc.vols)
			assert.EqualValues(t, tc.expected, actual)
		})
	}
}
