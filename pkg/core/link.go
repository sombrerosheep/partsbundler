package core

import (
	"fmt"
)

type Link struct {
	ID  int64  `json:"id"`
	URL string `json:"url"`
}

type LinkNotFound struct {
	LinkID, OwnerID int64
}

func (l LinkNotFound) Error() string {
	return fmt.Sprintf("Link %d not found on Part %d", l.LinkID, l.OwnerID)
}
