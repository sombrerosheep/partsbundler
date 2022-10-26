package main

type Storage interface {
	Connect() error
	Close() error

	// Links
	GetLinksForPart(int64) ([]Link, error)
	GetLinksForKit(int64) ([]Link, error)
	AddLinkToPart(int64, string) (Link, error)
	AddLinkToKit(int64, string) (Link, error)
	RemoveLinkFromPart(int64, int64) error
	RemoveLinkFromKit(int64, int64) error

	// Parts
	GetParts() ([]Part, error)
	GetPart(int64) (Part, error)
	AddPart(Part) (Part, error)
	UpdatePart(Part) (Part, error)
	DeletePart(int64) error

	// Kits
	GetKits() ([]Kit, error)
	GetKit(int64) (Kit, error)
	AddKit(Kit) (Kit, error)
	UpdateKit(Kit) (Kit, error)
	DeleteKit(int64) error
	GetKitParts(int64) ([]KitPart, error)
	AddPartToKit(int64, int64, uint64) error
	SetPartQuantityForKit(int64, uint64, int64) error
	RemovePartFromKit(int64, int64) error
}
