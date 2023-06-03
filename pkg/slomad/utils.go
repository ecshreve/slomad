package slomad

import (
	"bytes"
	"encoding/gob"

	nomadApi "github.com/hashicorp/nomad/api"
	nomadStructs "github.com/hashicorp/nomad/nomad/structs"
)

// StringPtr returns a pointer to the given string.
func StringPtr(s string) *string {
	return &s
}

// StringValOr returns the value of a string pointer if it's not nil, or a
// default value otherwise.
func StringValOr(sp *string, val string) string {
	if sp != nil {
		return *sp
	}
	return val
}

// IntPtr returns a pointer to the given int.
func IntPtr(i int) *int {
	return &i
}

// IntValOr returns the value of an int pointer if it's not nil, or a default.
func IntValOr(ip *int, val int) int {
	if ip != nil {
		return *ip
	}
	return val
}

// convertJob converts a Nomad Job to a Nomad API Job.
func convertJob(in *nomadStructs.Job) (*nomadApi.Job, error) {
	gob.Register([]map[string]interface{}{})
	gob.Register([]interface{}{})

	var apiJob *nomadApi.Job
	buf := new(bytes.Buffer)
	if err := gob.NewEncoder(buf).Encode(in); err != nil {
		return nil, err
	}
	if err := gob.NewDecoder(buf).Decode(&apiJob); err != nil {
		return nil, err
	}

	return apiJob, nil
}
