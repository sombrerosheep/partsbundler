package core

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

type Part struct {
	ID    int64    `json:"id"`
	Kind  PartType `json:"kind"`
	Name  string   `json:"name"`
	Links []Link   `json:"links"`
}
