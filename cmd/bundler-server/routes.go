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
		path:    "/kits",
		method:  http.MethodGet,
		handler: GetAllKits,
	},
	{
		path:    "/parts/:partId",
		method:  http.MethodGet,
		handler: GetPart,
	},
	{
		path:    "/kits/:kitId",
		method:  http.MethodGet,
		handler: GetKit,
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
