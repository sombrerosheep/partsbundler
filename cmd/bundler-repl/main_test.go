package main

import (
	"fmt"
	"strconv"
	"testing"

	"github.com/sombrerosheep/partsbundler/pkg/core"
	"github.com/stretchr/testify/assert"
)

func Test_GetCommand(t *testing.T) {
	tests := []struct {
		input    string
		expected ReplCmd
	}{
		{"get parts", GetPartsCmd{}},
		{"get part 1234", GetPartCmd{partId: 1234}},
		{"new part partType partName", NewPartCmd{name: "partName", kind: core.PartType("partType")}},
    {"delete part 1234", DeletePartCmd{partId: 1234}},
		{"add partlink 1234 example.com", AddPartLinkCmd{partId: 1234, link: "example.com"}},
    {"remove partlink 1234 789", RemovePartLinkCmd{partId: 1234, linkId: 789}},
		{"get kits", GetKitsCmd{}},
		{"get kit 1234", GetKitCmd{kitId: 1234}},
    {"new kit kitName kitSchem kitDiag", NewKitCmd{name: "kitName", schematic: "kitSchem", diagram: "kitDiag"}},
		{"delete kit 123", DeleteKitCmd{kitId: 123}},
    {"add kitlink 1234 example.com/kitlink", AddKitLinkCmd{kitId: 1234, link: "example.com/kitlink"}},
    {"remove kitlink 1234 789", RemoveKitLinkCmd{kitId: 1234, linkId: 789}},
    {"add kitpart 123 789 9", AddKitPartCmd{kitId: 123, partId: 789, quantity: 9}},
    {"set kitpart 123 789 5", SetKitPartQuantityCmd{kitId: 123, partId: 789, quantity: 5}},
    {"remove kitpart 123 789", RemoveKitPartCmd{kitId: 123, partId: 789}},
	}

	for _, test := range tests {
		t.Run(fmt.Sprintf("should return %T", test.expected), func(t *testing.T) {
			cmd, err := GetCommand(test.input)

			assert.Nil(t, err)
			assert.Equal(t, test.expected, cmd)
		})
	}
}

func Test_GetCommand_Errors(t *testing.T) {
  tests := []struct{
    input string
    errType error
    }{
      {"get part abd", &strconv.NumError{}},
      {"get part", CannotParseCommand{}},
      {"get kit", CannotParseCommand{}},
      {"get kit abc", &strconv.NumError{}},
  }

  for _, test := range tests {
    t.Run(fmt.Sprintf("Input (%s) should produce %T error", test.input, test.errType), func(t *testing.T) {
      cmd, err := GetCommand(test.input)

      assert.Nil(t, cmd)
      assert.IsType(t, test.errType, err)
    })
  }
}
