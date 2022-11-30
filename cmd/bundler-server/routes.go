package main

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/sombrerosheep/partsbundler/pkg/core"
)

type Endpoint struct {
	path    string
	method  string
	handler gin.HandlerFunc
}

var endpoints = []Endpoint{
	{
		path:    "/parts",
		method:  http.MethodGet,
		handler: GetAllParts,
	},
	{
		path:    "/parts/:partId",
		method:  http.MethodGet,
		handler: GetPart,
	},
	{
		path:    "/parts",
		method:  http.MethodPost,
		handler: CreatePart,
	},
	{
		path:    "/parts/:partId",
		method:  http.MethodDelete,
		handler: DeletePart,
	},
	{
		path:    "/parts/:partId/links",
		method:  http.MethodPost,
		handler: AddPartLink,
	},
	{
		path:    "/parts/:partId/links/:linkId",
		method:  http.MethodDelete,
		handler: RemovePartLink,
	},
	{
		path:    "/kits",
		method:  http.MethodGet,
		handler: GetAllKits,
	},
	{
		path:    "/kits/:kitId",
		method:  http.MethodGet,
		handler: GetKit,
	},
	{
		path:    "/kits",
		method:  http.MethodPost,
		handler: CreateKit,
	},
	{
		path:    "/kits/:kitId",
		method:  http.MethodDelete,
		handler: DeleteKit,
	},
	{
		path:    "/kits/:kitId/links",
		method:  http.MethodPost,
		handler: AddKitLink,
	},
	{
		path:    "/kits/:kitId/links/:linkId",
		method:  http.MethodDelete,
		handler: RemoveKitLink,
	},
	{
		path:    "/kits/:kitId/parts/:partId",
		method:  http.MethodPost,
		handler: AddKitPart,
	},
	{
		path:    "/kits/:kitId/parts/:partId",
		method:  http.MethodDelete,
		handler: RemoveKitPart,
	},
	{
		path:    "/kits/:kitId/parts/:partId/:quantity",
		method:  http.MethodPut,
		handler: UpdateKitPartQuantity,
	},
}

func GetAllParts(c *gin.Context) {
	svc := GetBundlerService()
	parts, err := svc.Parts.GetAll()
	if err != nil {
		c.String(http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, parts)
}

func GetPart(c *gin.Context) {
	svc := GetBundlerService()
	sid := c.Param("partId")
	id, err := strconv.ParseInt(sid, 10, 64)
	if err != nil {
		c.String(http.StatusBadRequest, err.Error())
		return
	}

	part, err := svc.Parts.Get(id)
	if err != nil {
		if _, ok := err.(core.PartNotFound); ok {
			c.String(http.StatusNotFound, err.Error())
			return
		}

		c.String(http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, part)
}

func CreatePart(c *gin.Context) {
	svc := GetBundlerService()

	var input core.Part
	err := c.BindJSON(&input)
	if err != nil {
		c.String(http.StatusBadRequest, err.Error())
		return
	}

	part, err := svc.Parts.New(input.Name, core.PartType(input.Kind))
	if err != nil {
		c.String(http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, part)
}

func DeletePart(c *gin.Context) {
	svc := GetBundlerService()
	partId := c.Param("partId")

	id, err := strconv.ParseInt(partId, 10, 64)
	if err != nil {
		c.String(http.StatusBadRequest, err.Error())
		return
	}

	err = svc.Parts.Delete(id)
	if err != nil {
		if _, ok := err.(core.PartNotFound); ok {
			c.String(http.StatusNotFound, err.Error())
			return
		}

		c.String(http.StatusInternalServerError, err.Error())
		return
	}

	c.Status(http.StatusNoContent)
}

func GetAllKits(c *gin.Context) {
	svc := GetBundlerService()
	kits, err := svc.Kits.GetAll()
	if err != nil {
		c.String(http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, kits)
}

func GetKit(c *gin.Context) {
	svc := GetBundlerService()

	sid := c.Param("kitId")
	id, err := strconv.ParseInt(sid, 10, 64)
	if err != nil {
		c.String(http.StatusBadRequest, err.Error())
		return
	}

	kit, err := svc.Kits.Get(id)
	if err != nil {
		if _, ok := err.(core.KitNotFound); ok {
			c.String(http.StatusNotFound, err.Error())
			return
		}

		c.String(http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, kit)
}

func AddPartLink(c *gin.Context) {
	svc := GetBundlerService()

	sid := c.Param("partId")
	id, err := strconv.ParseInt(sid, 10, 64)
	if err != nil {
		c.String(http.StatusBadRequest, err.Error())
		return
	}

	_, err = svc.Parts.Get(id)
	if err != nil {
		if _, ok := err.(core.PartNotFound); ok {
			c.String(http.StatusNotFound, err.Error())
			return
		}

		c.String(http.StatusInternalServerError, err.Error())
		return
	}

	var link core.Link
	err = c.BindJSON(&link)
	if err != nil {
		c.String(http.StatusBadRequest, err.Error())
		return
	}

	newLink, err := svc.Parts.AddLink(id, link.URL)
	if err != nil {
		c.String(http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, newLink)
}

func RemovePartLink(c *gin.Context) {
	svc := GetBundlerService()

	sid := c.Param("partId")
	partId, err := strconv.ParseInt(sid, 10, 64)
	if err != nil {
		c.String(http.StatusBadRequest, err.Error())
		return
	}

	sid = c.Param("linkId")
	linkId, err := strconv.ParseInt(sid, 10, 64)
	if err != nil {
		c.String(http.StatusBadRequest, err.Error())
		return
	}

	if linkId < 1 {
		c.String(http.StatusBadRequest, "linkId must not be valide")
		return
	}

	err = svc.Parts.RemoveLink(partId, linkId)
	if err != nil {
		c.String(http.StatusInternalServerError, err.Error())
		return
	}

	c.Status(http.StatusNoContent)

}

func CreateKit(c *gin.Context) {
	svc := GetBundlerService()

	var input core.Kit
	c.BindJSON(&input)

	kit, err := svc.Kits.New(input.Name, input.Schematic, input.Diagram)
	if err != nil {
		c.String(http.StatusInternalServerError, err.Error())
	}

	c.JSON(http.StatusOK, kit)
}

func DeleteKit(c *gin.Context) {
	svc := GetBundlerService()
	kitId := c.Param("kitId")

	id, err := strconv.ParseInt(kitId, 10, 64)
	if err != nil {
		c.String(http.StatusBadRequest, err.Error())
		return
	}

	err = svc.Kits.Delete(id)
	if err != nil {
		if _, ok := err.(core.KitNotFound); ok {
			c.String(http.StatusNotFound, err.Error())
			return
		}
		c.String(http.StatusInternalServerError, err.Error())
		return
	}

	c.Status(http.StatusNoContent)
}

func AddKitLink(c *gin.Context) {
	svc := GetBundlerService()
	kitId := c.Param("kitId")

	id, err := strconv.ParseInt(kitId, 10, 64)
	if err != nil {
		c.String(http.StatusBadRequest, err.Error())
		return
	}

	_, err = svc.Kits.Get(id)
	if err != nil {
		if _, ok := err.(core.KitNotFound); ok {
			c.String(http.StatusNotFound, err.Error())
			return
		}

		c.String(http.StatusInternalServerError, err.Error())
		return
	}

	var input core.Link
	err = c.BindJSON(&input)
	if err != nil {
		c.String(http.StatusBadRequest, err.Error())
		return
	}

	link, err := svc.Kits.AddLink(id, input.URL)
	if err != nil {
		c.String(http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, link)
}

func RemoveKitLink(c *gin.Context) {
	svc := GetBundlerService()

	sid := c.Param("kitId")
	kitId, err := strconv.ParseInt(sid, 10, 64)
	if err != nil {
		c.String(http.StatusBadRequest, err.Error())
		return
	}

	sid = c.Param("linkId")
	linkId, err := strconv.ParseInt(sid, 10, 64)
	if err != nil {
		c.String(http.StatusBadRequest, err.Error())
		return
	}

	err = svc.Kits.RemoveLink(kitId, linkId)
	if err != nil {
		c.String(http.StatusInternalServerError, err.Error())
		return
	}

	c.Status(http.StatusNoContent)
}

func AddKitPart(c *gin.Context) {
	svc := GetBundlerService()

	defaultQuantity := uint64(1)

	sid := c.Param("kitId")
	kitId, err := strconv.ParseInt(sid, 10, 64)
	if err != nil {
		c.String(http.StatusBadRequest, err.Error())
		return
	}

	sid = c.Param("partId")
	partId, err := strconv.ParseInt(sid, 10, 64)
	if err != nil {
		c.String(http.StatusBadRequest, err.Error())
		return
	}

	sqty := c.Query("quantity")
	qty, err := strconv.ParseUint(sqty, 10, 64)
	if err != nil {
		qty = defaultQuantity
	}

	err = svc.Kits.AddPart(kitId, partId, qty)
	if err != nil {
		c.String(http.StatusInternalServerError, err.Error())
		return
	}

	part, err := svc.Parts.Get(partId)
	if err != nil {
		c.String(http.StatusInternalServerError, err.Error())
		return
	}

	kitPart := core.KitPart{
		Part:     part,
		Quantity: qty,
	}

	c.JSON(http.StatusOK, kitPart)
}

// todo: Remove Kit Part
func RemoveKitPart(c *gin.Context) {
	svc := GetBundlerService()

	sid := c.Param("kitId")
	kitId, err := strconv.ParseInt(sid, 10, 64)
	if err != nil {
		c.String(http.StatusBadRequest, err.Error())
		return
	}

	sid = c.Param("partId")
	partId, err := strconv.ParseInt(sid, 10, 64)
	if err != nil {
		c.String(http.StatusBadRequest, err.Error())
		return
	}

	err = svc.Kits.RemovePart(kitId, partId)
	if err != nil {
		c.String(http.StatusInternalServerError, err.Error())
		return
	}

	c.Status(http.StatusNoContent)
}

// todo: Update Kit Part Quantity
func UpdateKitPartQuantity(c *gin.Context) {
	svc := GetBundlerService()

	sid := c.Param("kitId")
	kitId, err := strconv.ParseInt(sid, 10, 64)
	if err != nil {
		c.String(http.StatusBadRequest, err.Error())
		return
	}

	sid = c.Param("partId")
	partId, err := strconv.ParseInt(sid, 10, 64)
	if err != nil {
		c.String(http.StatusBadRequest, err.Error())
		return
	}

	sid = c.Param("quantity")
	quantity, err := strconv.ParseUint(sid, 10, 64)
	if err != nil {
		c.String(http.StatusBadRequest, err.Error())
		return
	}

	err = svc.Kits.SetPartQuantity(kitId, partId, quantity)
	if err != nil {
		c.String(http.StatusInternalServerError, err.Error())
		return
	}

	part, err := svc.Parts.Get(partId)
	if err != nil {
		c.String(http.StatusInternalServerError, err.Error())
		return
	}

	kitPart := core.KitPart{
		Part:     part,
		Quantity: quantity,
	}

	c.JSON(http.StatusOK, kitPart)

}
