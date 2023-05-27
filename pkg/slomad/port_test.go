package slomad_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/ecshreve/slomad/pkg/slomad"
	nomadStructs "github.com/hashicorp/nomad/nomad/structs"
)

func TestToNomadPort(t *testing.T) {
	tests := []struct {
		name     string
		args     slomad.Port
		expected nomadStructs.Port
	}{
		{
			name: "dynamic",
			args: *slomad.NewPort("http", 80, nil),
			expected: nomadStructs.Port{
				Label: "http",
				To:    80,
			},
		},
		{
			name: "static",
			args: *slomad.NewPort("http", 80, slomad.IntPtr(80)),
			expected: nomadStructs.Port{
				Label: "http",
				To:    80,
				Value: 80,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := tt.args
			assert.Equal(t, tt.expected, p.ToNomadPort())
		})
	}
}

func TestPortList(t *testing.T) {
	tests := []struct {
		name          string
		args          []slomad.Port
		expectDynamic []nomadStructs.Port
		expectStatic  []nomadStructs.Port
	}{
		{
			name:          "empty",
			args:          []slomad.Port{},
			expectDynamic: []nomadStructs.Port{},
			expectStatic:  []nomadStructs.Port{},
		},
		{
			name: "one dynamic",
			args: []slomad.Port{
				*slomad.NewPort("http", 80, nil),
			},
			expectDynamic: []nomadStructs.Port{
				{
					Label: "http",
					To:    80,
				},
			},
			expectStatic: []nomadStructs.Port{},
		},
		{
			name: "one static",
			args: []slomad.Port{
				*slomad.NewPort("http", 80, slomad.IntPtr(80)),
			},
			expectDynamic: []nomadStructs.Port{},
			expectStatic: []nomadStructs.Port{
				{
					Label: "http",
					To:    80,
					Value: 80,
				},
			},
		},
		{
			name: "one dynamic, one static",
			args: []slomad.Port{
				*slomad.NewPort("http", 80, nil),
				*slomad.NewPort("https", 443, slomad.IntPtr(443)),
			},
			expectDynamic: []nomadStructs.Port{
				{
					Label: "http",
					To:    80,
				},
			},
			expectStatic: []nomadStructs.Port{
				{
					Label: "https",
					To:    443,
					Value: 443,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			pl := slomad.NewPortList(tt.args)
			assert.Equal(t, tt.expectDynamic, pl.GetDynamic())
			assert.Equal(t, tt.expectStatic, pl.GetStatic())
		})
	}
}

func TestBuildDynamicPortList(t *testing.T) {
	tests := []struct {
		name          string
		args          []int
		labelBase     string
		expectDynamic []nomadStructs.Port
	}{
		{
			name:          "empty",
			args:          []int{},
			labelBase:     "http",
			expectDynamic: []nomadStructs.Port{},
		},
		{
			name:      "one",
			args:      []int{80},
			labelBase: "http",
			expectDynamic: []nomadStructs.Port{
				{
					Label: "http-0",
					To:    80,
				},
			},
		},
		{
			name:      "two",
			args:      []int{80, 443},
			labelBase: "http",
			expectDynamic: []nomadStructs.Port{
				{
					Label: "http-0",
					To:    80,
				},
				{
					Label: "http-1",
					To:    443,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			pl := slomad.BuildDynamicPortList(tt.args)
			assert.Equal(t, tt.expectDynamic, pl.GetDynamic())
		})
	}
}
