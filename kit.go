package main

type Kit struct {
	ID        int64     `json:"id"`
	Parts     []KitPart `json:"parts"`
	Name      string    `json:"name"`
	Schematic string    `json:"schematics"`
	Diagram   string    `json:"diagram,omitempty"`
	Links     []string  `json:"links,omitempty"`
}

type KitPart struct {
	Part
	Quantity uint64 `json:"quantity"`
}
