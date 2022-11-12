package core

import (
	"fmt"
)

type InvalidPartType struct {
	InvalidType string
}

func (p InvalidPartType) Error() string {
	return fmt.Sprintf("Invalid PartType '%s'", p.InvalidType)
}

type PartType string

const (
	Resistor      PartType = "Resistor"
	Capacitor              = "Capacitor"
	IC                     = "IC"
	Transistor             = "Transistor"
	Diode                  = "Diode"
	Potentiometer          = "Potentiometer"
	Switch                 = "Switch"
)

func (p PartType) IsValid() error {
	switch p {
	case Resistor, Capacitor, IC, Transistor, Diode, Potentiometer, Switch:
		return nil
	}

	return InvalidPartType{string(p)}
}

type Part struct {
	ID    int64    `json:"id"`
	Kind  PartType `json:"kind"`
	Name  string   `json:"name"`
	Links []Link   `json:"links"`
}

type PartNotFound struct {
	PartID int64
}

func (p PartNotFound) Error() string {
	return fmt.Sprintf("Part %d not found", p.PartID)
}
