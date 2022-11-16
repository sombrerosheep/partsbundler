package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/sombrerosheep/partsbundler/pkg/core"
	"github.com/sombrerosheep/partsbundler/pkg/service/mock"
	"github.com/stretchr/testify/assert"
)

func CreateStubServer() *gin.Engine {
	router := gin.Default()

	RegisterEndpoints(router, endpoints)

	return router
}

func Test_GetAllParts(t *testing.T) {
	t.Run("should return parts", func(t *testing.T) {
		router := CreateStubServer()
		bundlerService = mock.StubBundlerService
		w := httptest.NewRecorder()
		req, err := http.NewRequest(http.MethodGet, "/parts", nil)
		assert.Nil(t, err)

		router.ServeHTTP(w, req)

		var actualParts []core.Part

		err = json.Unmarshal(w.Body.Bytes(), &actualParts)

		assert.Nil(t, err)
		assert.Equal(t, http.StatusOK, w.Code)
		assert.Equal(t, mock.FakeParts[:], actualParts)
	})
}

func Test_GetPart(t *testing.T) {
	t.Run("should get each part", func(t *testing.T) {
		router := CreateStubServer()
		bundlerService = mock.StubBundlerService

		for _, v := range mock.FakeParts {

			w := httptest.NewRecorder()
			req, err := http.NewRequest(http.MethodGet, fmt.Sprintf("/parts/%d", v.ID), nil)
			assert.Nil(t, err)

			router.ServeHTTP(w, req)

			var part core.Part

			err = json.Unmarshal(w.Body.Bytes(), &part)

			assert.Nil(t, err)
			assert.Equal(t, http.StatusOK, w.Code)
			assert.Equal(t, v, part)
		}
	})

	t.Run("should return bad request if partId is invalid", func(t *testing.T) {
		router := CreateStubServer()
		bundlerService = mock.StubBundlerService

		w := httptest.NewRecorder()
		req, err := http.NewRequest(http.MethodGet, "/parts/onetwothree", nil)

		assert.Nil(t, err)

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("should return PartNotFound when kit does not exist", func(t *testing.T) {
		router := CreateStubServer()
		bundlerService = mock.StubBundlerService

		partId := int64(9999)

		w := httptest.NewRecorder()
		req, err := http.NewRequest(http.MethodGet, fmt.Sprintf("/parts/%d", partId), nil)

		assert.Nil(t, err)

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusNotFound, w.Code)
		assert.Equal(t, core.PartNotFound{PartID: partId}.Error(), w.Body.String())
	})
}

func Test_GetAllKits(t *testing.T) {
	t.Run("should return kits", func(t *testing.T) {
		router := CreateStubServer()
		bundlerService = mock.StubBundlerService

		w := httptest.NewRecorder()
		req, err := http.NewRequest(http.MethodGet, "/kits", nil)

		assert.Nil(t, err)

		router.ServeHTTP(w, req)

		var actualKits []core.Kit

		err = json.Unmarshal(w.Body.Bytes(), &actualKits)

		assert.Nil(t, err)
		assert.Equal(t, http.StatusOK, w.Code)
		assert.Equal(t, mock.FakeKits[:], actualKits)
	})
}

func Test_GetKit(t *testing.T) {
	t.Run("should each kit", func(t *testing.T) {
		router := CreateStubServer()
		bundlerService = mock.StubBundlerService

		for _, v := range mock.FakeKits {

			w := httptest.NewRecorder()
			req, err := http.NewRequest(http.MethodGet, fmt.Sprintf("/kits/%d", v.ID), nil)

			assert.Nil(t, err)

			router.ServeHTTP(w, req)

			var kit core.Kit

			err = json.Unmarshal(w.Body.Bytes(), &kit)

			assert.Nil(t, err)
			assert.Equal(t, http.StatusOK, w.Code)
			assert.Equal(t, v, kit)
		}
	})

	t.Run("should return bad request if kitId is invalid", func(t *testing.T) {
		router := CreateStubServer()
		bundlerService = mock.StubBundlerService

		w := httptest.NewRecorder()
		req, err := http.NewRequest(http.MethodGet, "/kits/onetwothree", nil)

		assert.Nil(t, err)

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("should return KitNotFound when kit does not exist", func(t *testing.T) {
		router := CreateStubServer()
		bundlerService = mock.StubBundlerService

		kitId := int64(9999)

		w := httptest.NewRecorder()
		req, err := http.NewRequest(http.MethodGet, fmt.Sprintf("/kits/%d", kitId), nil)

		assert.Nil(t, err)

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusNotFound, w.Code)
		assert.Equal(t, core.KitNotFound{KitID: kitId}.Error(), w.Body.String())

	})
}

func Test_CreateKit(t *testing.T) {
	t.Run("should return created kit", func(t *testing.T) {
		router := CreateStubServer()
		bundlerService = mock.StubBundlerService

		kitName := "my kit"
		kitSchem := "example.com/my-schematic"
		kitDiag := "example.com/my-diag"

		input := KitCreateInput{
			Name:      kitName,
			Schematic: kitSchem,
			Diagram:   kitDiag,
		}
		inBytes, err := json.Marshal(input)

		assert.Nil(t, err)

		reqBody := bytes.NewReader(inBytes)

		w := httptest.NewRecorder()
		req, err := http.NewRequest(http.MethodPost, "/kits", reqBody)

		assert.Nil(t, err)

		router.ServeHTTP(w, req)

		fmt.Println(string(w.Body.Bytes()))

		var kit core.Kit
		err = json.Unmarshal(w.Body.Bytes(), &kit)

		assert.Nil(t, err)

		assert.Equal(t, kitName, kit.Name)
		assert.Equal(t, kitSchem, kit.Schematic)
		assert.Equal(t, kitDiag, kit.Diagram)
	})
}
