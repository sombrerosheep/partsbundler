package main

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type Endpoint struct {
	path    string
	method  string
	handler gin.HandlerFunc
}

var endpoints = []Endpoint {
	{
		path:    "/parts",
		method:  "GET",
		handler: GetAllParts,
	},
	{
    path:    "/kits",
		method:  "GET",
		handler: GetAllKits,
	},
  {
    path: "/parts/:partId",
    method: "GET",
    handler: GetPart,
  },
}

func GetAllParts(c *gin.Context) {
	svc := GetBundlerService()
	parts, err := svc.Parts.GetAll()
	if err != nil {
		c.String(http.StatusInternalServerError, err.Error())
	}

	c.JSON(http.StatusOK, parts)
}

func GetPart(c *gin.Context) {
	svc := GetBundlerService()
	sid := c.Param("partId")
	id, err := strconv.ParseInt(sid, 10, 64)
	if err != nil {
		c.String(http.StatusBadRequest, err.Error())
	}

	part, err := svc.Parts.Get(id)
	if err != nil {
		c.String(http.StatusInternalServerError, err.Error())
	}

	c.JSON(http.StatusOK, part)
}

func GetAllKits(c *gin.Context) {
	svc := GetBundlerService()
	kits, err := svc.Kits.GetAll()
	if err != nil {
		c.String(http.StatusInternalServerError, err.Error())
	}

	c.JSON(http.StatusOK, kits)
}
