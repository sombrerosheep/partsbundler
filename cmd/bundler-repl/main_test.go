package main

import (
	"fmt"
	"testing"

	"github.com/sombrerosheep/partsbundler/pkg/core"
	"github.com/stretchr/testify/assert"
)

func Test_GetCommand(t *testing.T) {
	tests := []struct {
		input    string
		expected ReplCmd
		err      error
	}{
		{"get parts", GetPartsCmd{}, nil},
		{"get kits", GetKitsCmd{}, nil},
		{"get part 1234", GetPartCmd{partId: 1234}, nil},
		{"get kit 1234", GetKitCmd{kitId: 1234}, nil},
		{"new part partType partName", NewPartCmd{name: "partName", kind: core.PartType("partType")}, nil},
		{"add partlink 1234 example.com", AddPartLinkCommand{partId: 1234, link: "example.com"}, nil},
		{"delete part 123", DeletePartCommand{partId: 123}, nil},
	}

	for _, test := range tests {
		t.Run(fmt.Sprintf("should return %T", test.expected), func(t *testing.T) {
			cmd, err := GetCommand(test.input)

			assert.Equal(t, test.err, err)
			assert.Equal(t, test.expected, cmd)
		})
	}
}
