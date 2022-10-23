package main

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
	ID    uint64   `json:"id"`
	Kind  PartType `json:"kind"`
	Name  string   `json:"name"`
	Links []string `json:"links"`
}

func (d *SqlLiteDb) GetPart(partID uint64) (Part, error) {
	const stmt string = `
		SELECT * from 
	`

	return Part{}, nil
}

func (d *SqlLiteDb) PutPart(p Part) (Part, error) {
	return p, nil
}

func (d *SqlLiteDb) UpdatePart(p Part) (Part, error) {
	return p, nil
}

func (d *SqlLiteDb) DeletePart(partID uint64) error {
	return nil
}