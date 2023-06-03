package slomad_test

import (
	"testing"

	"github.com/ecshreve/slomad/pkg/slomad"
	"github.com/stretchr/testify/assert"
)

func TestNewPort(t *testing.T) {
	p := slomad.NewPort(slomad.PortParams{
		Label: "http",
		To:    8080,
		From:  8080,
	})

	assert.Equal(t, p.Label, "http")
	assert.Equal(t, p.To, 8080)
	assert.Equal(t, p.From, 8080)
	assert.Equal(t, p.Static, true)
}

func TestNewPorts(t *testing.T) {
	ports := slomad.NewPorts([]slomad.PortParams{
		{Label: "http", To: 8080, From: 8080},
		{Label: "https", To: 8443},
	})

	assert.Equal(t, len(ports), 2)

	assert.Equal(t, ports[0].Label, "http")
	assert.Equal(t, ports[0].To, 8080)
	assert.Equal(t, ports[0].From, 8080)
	assert.Equal(t, ports[0].Static, true)

	assert.Equal(t, ports[1].Label, "https")
	assert.Equal(t, ports[1].To, 8443)
	assert.Equal(t, ports[1].From, 0)
	assert.Equal(t, ports[1].Static, false)
}
