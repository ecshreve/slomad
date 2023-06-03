package slomad

import (
	"testing"

	nomadStructs "github.com/hashicorp/nomad/nomad/structs"
	"github.com/samsarahq/go/snapshotter"
	"github.com/stretchr/testify/assert"
)

func TestToNomadPort(t *testing.T) {
	testcases := []struct {
		desc     string
		port     Port
		expected nomadStructs.Port
	}{
		{
			desc:     "nil port",
			port:     Port{},
			expected: nomadStructs.Port{},
		},
		{
			desc:     "empty port",
			port:     Port{Label: "http"},
			expected: nomadStructs.Port{Label: "http", To: 0, Value: 0},
		},
		{
			desc:     "dynamic port",
			port:     Port{Label: "http", To: 8080, From: 0, Static: false},
			expected: nomadStructs.Port{Label: "http", To: 8080, Value: 0},
		},
		{
			desc:     "static port",
			port:     Port{Label: "ssh", To: 22, From: 22, Static: true},
			expected: nomadStructs.Port{Label: "ssh", To: 22, Value: 22},
		},
	}

	for _, tc := range testcases {
		t.Run(tc.desc, func(t *testing.T) {
			assert.Equal(t, tc.expected, toNomadPort(&tc.port))
		})
	}
}

func TestBasicPort(t *testing.T) {
	testcases := []struct {
		desc     string
		val      int
		expected Port
	}{
		{
			desc:     "zero value port",
			val:      0,
			expected: Port{Label: "http", To: 0, From: 0, Static: false},
		},
		{
			desc:     "normal value",
			val:      88,
			expected: Port{Label: "http", To: 88, From: 0, Static: false},
		},
	}

	for _, tc := range testcases {
		t.Run(tc.desc, func(t *testing.T) {
			assert.Equal(t, tc.expected, *basicPort(tc.val))
			assert.Equal(t, []*Port{&tc.expected}, BasicPortConfig(tc.val))
		})
	}
}

func TestExtractLabels(t *testing.T) {
	ports := []*Port{
		{Label: "http", To: 8080, From: 0, Static: false},
		{Label: "ssh", To: 22, From: 22, Static: true},
		{Label: "https", To: 4334, From: 0, Static: false},
		{Label: "dns", To: 8081, From: 80, Static: true},
	}

	expected := []string{"http", "ssh", "https", "dns"}
	actual := extractLabels(ports)
	assert.Equal(t, expected, actual)
}

func TestToNomadPortMap(t *testing.T) {
	ports := []*Port{
		{Label: "http", To: 8080, From: 0, Static: false},
		{Label: "https", To: 4334, From: 0, Static: false},
		{Label: "dns", To: 8081, From: 80, Static: true},
		{Label: "ssh", To: 22, From: 22, Static: true},
	}

	nomadPorts := []nomadStructs.Port{
		{Label: "http", To: 8080, Value: 0},
		{Label: "https", To: 4334, Value: 0},
		{Label: "dns", To: 8081, Value: 80},
		{Label: "ssh", To: 22, Value: 22},
	}

	expected := map[string][]nomadStructs.Port{
		"dynamic": {nomadPorts[0], nomadPorts[1]},
		"static":  {nomadPorts[2], nomadPorts[3]},
	}

	actual := toNomadPortMap(ports)
	assert.Equal(t, expected["dynamic"], actual["dynamic"])
	assert.Equal(t, expected["static"], actual["static"])
}

func TestGetNetworks(t *testing.T) {
	snap := snapshotter.New(t)
	defer snap.Verify()

	ports := []*Port{
		{Label: "http", To: 8080, From: 0, Static: false},
		{Label: "https", To: 4334, From: 0, Static: false},
		{Label: "dns", To: 8081, From: 80, Static: true},
		{Label: "ssh", To: 22, From: 22, Static: true},
	}

	actual := getNetworks(ports)
	snap.Snapshot("networks", actual)
}
