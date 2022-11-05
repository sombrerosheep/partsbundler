package main

type IPartService interface {
	GetAll() ([]Part, error)
	Get(partId int64) (Part, error)

	AddLink(partId int64, link string) (Link, error)
	RemoveLink(partId int64, linkId int64) error

	New(name string, kind PartType) (Part, error)
	Delete(partId int64) error
}

type IKitService interface {
	GetAll() ([]Kit, error)
	Get(kitId int64) (Kit, error)

	AddLink(kitId int64, link string) (Link, error)
	RemoveLink(kitId int64, linkId int64) error

	AddPart(kitId int64, partId int64, quantity uint64) error
	SetPartQuantity(kitId int64, partId int64, quantity uint64) error
	RemovePart(kitId int64, partId int64) error

	New(name string, schematic string, diagram string) (Kit, error)
	Delete(kitId int64) error
}
