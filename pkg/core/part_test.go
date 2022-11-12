package core

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_PartType_IsValid(t *testing.T) {
  tests := []struct{
    input string
    expectValid bool
  }{
    {"Resistor", true},
    {"Capacitor", true},
    {"IC", true},
    {"Transistor", true},
    {"Diode", true},
    {"Potentiometer", true},
    {"Switch", true},
    {"", false},
    {"Flux Capacitor", false},
  }

  for _, test := range tests {
    t.Run(fmt.Sprintf("Expect '%s' to be %t", test.input, test.expectValid), func(t *testing.T) {
      pt := PartType(test.input)

      err := pt.IsValid()

      if test.expectValid == true {
        assert.Nil(t, err)
      } else {
        assert.NotNil(t, err)
        assert.IsType(t, InvalidPartType{}, err)
        assert.Equal(t, test.input, err.(InvalidPartType).InvalidType)
      }
    })
  }
}
