package main

type Storage interface {
	Connect() error
	Close() error

	// Links
	GetLinksForPart(partId int64) ([]Link, error)
	GetLinksForKit(kitId int64) ([]Link, error)
	AddLinkToPart(partId int64, url string) (Link, error)
	AddLinkToKit(kitId int64, url string) (Link, error)
	RemoveLinkFromPart(partId int64, linkId int64) error
	RemoveLinkFromKit(kitId int64, linkId int64) error

	// Parts
	GetParts() ([]Part, error)
	GetPart(partId int64) (Part, error)
	AddPart(Part) (Part, error)
	UpdatePart(Part) (Part, error)
	DeletePart(partId int64) error

	// Kits
	GetKits() ([]Kit, error)
	GetKit(kitId int64) (Kit, error)
	AddKit(Kit) (Kit, error)
	UpdateKit(Kit) (Kit, error)
	DeleteKit(kitId int64) error
	GetKitParts(kitId int64) ([]KitPart, error)
	AddPartToKit(partId int64, kitId int64, quantity uint64) error
	SetPartQuantityForKit(partId int64, kitId uint64, quantity int64) error
	RemovePartFromKit(partId int64, kitId int64) error
}
