package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_generator(t *testing.T) {
	tests := []struct {
		name   string
		fields []Field
		expect string
	}{
		// TODO: Add test cases.
		{
			fields: []Field{
				Field{
					Name: "Name",
					Type: "string",
				},
				Field{
					Name: "Value",
					Type: "int",
				},
			},
			expect: `map[string]int`,
		},
		{
			fields: []Field{
				Field{
					Name: "AppID",
					Type: "string",
				},
				Field{
					Name: "Cluster",
					Type: "string",
				},
				Field{
					Name: "IP",
					Type: "string",
				},
			},
			expect: `map[string]map[string]string`,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.True(t, tt.expect == generateType(tt.fields))
		})
	}
}
