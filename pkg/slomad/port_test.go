package slomad_test

import (
	"testing"

	"github.com/ecshreve/slomad/pkg/slomad"
	"github.com/stretchr/testify/assert"
)

func TestNewPort(t *testing.T) {
	p := slomad.NewPort("http", 8080, 8080)

	assert.Equal(t, p.Label, "http")
	assert.Equal(t, p.To, 8080)
	assert.Equal(t, p.From, 8080)
	assert.Equal(t, p.Static, true)
}
