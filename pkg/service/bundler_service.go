package service

import (
	"github.com/sombrerosheep/partsbundler/pkg/core"
)

type IPartService interface {
	GetAll() ([]core.Part, error)
	Get(partId int64) (core.Part, error)

	AddLink(partId int64, link string) (core.Link, error)
	RemoveLink(partId int64, linkId int64) error

	New(name string, kind core.PartType) (core.Part, error)
	Delete(partId int64) error
}

type IKitService interface {
	GetAll() ([]core.Kit, error)
	Get(kitId int64) (core.Kit, error)

	AddLink(kitId int64, link string) (core.Link, error)
	RemoveLink(kitId int64, linkId int64) error

	AddPart(kitId int64, partId int64, quantity uint64) error
	GetPartUsage(partId int64) ([]int64, error)
	SetPartQuantity(kitId int64, partId int64, quantity uint64) error
	RemovePart(kitId int64, partId int64) error

	New(name string, schematic string, diagram string) (core.Kit, error)
	Delete(kitId int64) error
}

type BundlerService struct {
	Parts IPartService
	Kits  IKitService
}
