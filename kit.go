package main

type Kit struct {
	ID        int64    `json:"id"`
	Parts     []Part   `json:"parts"`
	Name      string   `json:"name"`
	Schematic string   `json:"schematics"`
	Diagram   string   `json:"diagram,omitempty"`
	Links     []string `json:"links,omitempty"`
}

func (d *SqlLiteDb) GetKit(kitID uint64) (Kit, error) {
	const stmt string = `
		SELECT * from 
	`

	return Kit{}, nil
}

// Kit Related
// create
// delete
// setSchematic
// setDiagram

// Links Related
// get links for kit
// add link
// remove link

// Parts Related

func (d *SqlLiteDb) GetKitParts(kitID uint64) ([]Part, error) {

	return []Part{}, nil
}

func (d *SqlLiteDb) AddPartToKit(p Part) (Part, error) {
	return p, nil
}

func (d *SqlLiteDb) RemovePartFromKit(p Part) (Part, error) {
	return p, nil
}
