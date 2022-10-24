package main

type Storage interface {
	Connect() error

	// Links
	GetLinksForPart(int64) ([]Link, error)
	GetLinksForKit(int64) ([]Link, error)
	AddLinkToPart(string) (Link, error)
	AddLinkToKit(string) (Link, error)
	RemoveLinkFromPart(int64, int64) error
	RemoveLinkFromKit(int64, int64) error

	// Parts
	GetPart(int64) (Part, error)
	PutPart(Part) (Part, error)
	UpdatePart(Part) (Part, error)
	DeletePart(int64) error

	// Kits
	GetKit(int64) (Kit, error)
	PutKit(Kit) (Kit, error)
	UpdateKit(Kit) (Kit, error)
	DeleteKit(int64) error
}