package core

import (
	"fmt"
)

type Kit struct {
	ID        int64     `json:"id"`
	Parts     []KitPart `json:"parts"`
	Name      string    `json:"name"`
	Schematic string    `json:"schematics"`
	Diagram   string    `json:"diagram,omitempty"`
	Links     []Link    `json:"links,omitempty"`
}

type KitPart struct {
	Part
	Quantity uint64 `json:"quantity"`
}

type KitNotFound struct {
	KitID int64
}

func (k KitNotFound) Error() string {
	return fmt.Sprintf("Kit %d not found", k.KitID)
}

type PartInUse struct {
	PartID int64
}

func (p PartInUse) Error() string {
	return fmt.Sprintf("Part %d is in use by one or more kits", p.PartID)
}
