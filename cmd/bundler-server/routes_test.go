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
	gin.SetMode(gin.TestMode)

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

func Test_CreatePart(t *testing.T) {
	t.Run("should create part", func(t *testing.T) {
		router := CreateStubServer()
		bundlerService = mock.StubBundlerService

		partName := "my part"
		partKind := "Capacitor"

		input := core.Part{
			Name: partName,
			Kind: core.PartType(partKind),
		}
		inBytes, err := json.Marshal(input)

		assert.Nil(t, err)

		reqBody := bytes.NewReader(inBytes)

		w := httptest.NewRecorder()
		req, err := http.NewRequest(http.MethodPost, "/parts", reqBody)

		assert.Nil(t, err)

		router.ServeHTTP(w, req)

		var part core.Part
		err = json.Unmarshal(w.Body.Bytes(), &part)

		assert.Nil(t, err)
		assert.Equal(t, partName, part.Name)
		assert.Equal(t, partKind, string(part.Kind))
	})
}

func Test_DeletePart(t *testing.T) {
	t.Run("should delete part", func(t *testing.T) {
		router := CreateStubServer()
		bundlerService = mock.StubBundlerService

		partId := int64(3)

		w := httptest.NewRecorder()
		req, err := http.NewRequest(http.MethodDelete, fmt.Sprintf("/parts/%d", partId), nil)

		assert.Nil(t, err)

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusNoContent, w.Result().StatusCode)
	})
}

func Test_AddPartLink(t *testing.T) {
	t.Run("should add link", func(t *testing.T) {
		router := CreateStubServer()
		bundlerService = mock.StubBundlerService

		partId := int64(1)
		newLink := core.Link{
			URL: "example.com/newlink",
		}

		buf, err := json.Marshal(newLink)

		assert.Nil(t, err)

		reader := bytes.NewReader(buf)

		w := httptest.NewRecorder()
		req, err := http.NewRequest(http.MethodPost, fmt.Sprintf("/parts/%d/links", partId), reader)

		assert.Nil(t, err)

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		var link core.Link
		err = json.Unmarshal(w.Body.Bytes(), &link)

		assert.Nil(t, err)
		assert.Greater(t, link.ID, int64(0))
		assert.Equal(t, newLink.URL, link.URL)
	})

	t.Run("should return bad request if body is invalid", func(t *testing.T) {
		router := CreateStubServer()
		bundlerService = mock.StubBundlerService

		partId := int64(1)

		reader := bytes.NewReader([]byte("{"))

		w := httptest.NewRecorder()
		req, err := http.NewRequest(http.MethodPost, fmt.Sprintf("/parts/%d/links", partId), reader)

		assert.Nil(t, err)

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("should return PartNotFound when part does not exist", func(t *testing.T) {
		router := CreateStubServer()
		bundlerService = mock.StubBundlerService

		partId := int64(999)

		reader := bytes.NewReader([]byte("{}"))

		w := httptest.NewRecorder()
		req, err := http.NewRequest(http.MethodPost, fmt.Sprintf("/parts/%d/links", partId), reader)

		assert.Nil(t, err)

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusNotFound, w.Code)
		assert.Equal(t, fmt.Sprintf("Part %d not found", partId), w.Body.String())
	})
}

func Test_RemovePartLink(t *testing.T) {
	t.Run("should remove part link", func(t *testing.T) {
		router := CreateStubServer()
		bundlerService = mock.StubBundlerService

		part := mock.FakeParts[0]
		link := part.Links[0]

		w := httptest.NewRecorder()
		uri := fmt.Sprintf("/parts/%d/links/%d", part.ID, link.ID)
		req, err := http.NewRequest(http.MethodDelete, uri, nil)

		assert.Nil(t, err)

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusNoContent, w.Code)
	})

	t.Run("should return PartNotFound when part does not exist", func(t *testing.T) {
		router := CreateStubServer()
		bundlerService = mock.StubBundlerService

		part := mock.FakeParts[0]
		link := part.Links[0]

		w := httptest.NewRecorder()
		uri := fmt.Sprintf("/parts/%d/links/%d", part.ID, link.ID)
		req, err := http.NewRequest(http.MethodDelete, uri, nil)

		assert.Nil(t, err)

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusNoContent, w.Code)
	})

	t.Run("should return LinkNotFound when link does not exist", func(t *testing.T) {
		router := CreateStubServer()
		bundlerService = mock.StubBundlerService

		part := mock.FakeParts[0]
		link := part.Links[0]

		w := httptest.NewRecorder()
		uri := fmt.Sprintf("/parts/%d/links/%d", part.ID, link.ID)
		req, err := http.NewRequest(http.MethodDelete, uri, nil)

		assert.Nil(t, err)

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusNoContent, w.Code)
	})
}

func Test_RemovePart(t *testing.T) {
	t.Run("should remove part", func(t *testing.T) {
		router := CreateStubServer()
		bundlerService = mock.StubBundlerService

		part := mock.FakeParts[0]

		w := httptest.NewRecorder()
		uri := fmt.Sprintf("/parts/%d", part.ID)
		req, err := http.NewRequest(http.MethodDelete, uri, nil)

		assert.Nil(t, err)

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusNoContent, w.Code)
	})
}

/////////////////////////
// Kit Tests

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

		input := core.Kit{
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

		var kit core.Kit
		err = json.Unmarshal(w.Body.Bytes(), &kit)

		assert.Nil(t, err)

		assert.Equal(t, kitName, kit.Name)
		assert.Equal(t, kitSchem, kit.Schematic)
		assert.Equal(t, kitDiag, kit.Diagram)
	})
}

func Test_DeleteKit(t *testing.T) {
	t.Run("should delete kit", func(t *testing.T) {
		router := CreateStubServer()
		bundlerService = mock.StubBundlerService

		kitId := int64(7777)

		w := httptest.NewRecorder()
		req, err := http.NewRequest(http.MethodDelete, fmt.Sprintf("/kits/%d", kitId), nil)

		assert.Nil(t, err)

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusNoContent, w.Result().StatusCode)
	})
}

func Test_AddKitLink(t *testing.T) {
	t.Run("should add link", func(t *testing.T) {
		router := CreateStubServer()
		bundlerService = mock.StubBundlerService

		kitId := mock.FakeKits[0].ID
		newLink := core.Link{
			URL: "example.com/newlink",
		}

		buf, err := json.Marshal(newLink)

		assert.Nil(t, err)

		reader := bytes.NewReader(buf)

		w := httptest.NewRecorder()
		req, err := http.NewRequest(http.MethodPost, fmt.Sprintf("/kits/%d/links", kitId), reader)

		router.ServeHTTP(w, req)

		assert.Nil(t, err)

		var link core.Link
		err = json.Unmarshal(w.Body.Bytes(), &link)

		assert.Nil(t, err)
		assert.Equal(t, http.StatusOK, w.Code)
		assert.Equal(t, newLink.URL, link.URL)
		assert.Greater(t, link.ID, int64(0))
	})

	t.Run("should return BadRequest if body is invalid", func(t *testing.T) {
		router := CreateStubServer()
		bundlerService = mock.StubBundlerService

		kitId := mock.FakeKits[0].ID

		reader := bytes.NewReader([]byte("{"))

		w := httptest.NewRecorder()
		req, err := http.NewRequest(http.MethodPost, fmt.Sprintf("/kits/%d/links", kitId), reader)

		router.ServeHTTP(w, req)

		assert.Nil(t, err)
		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("should return KitNotFound if kit does not exist", func(t *testing.T) {
		router := CreateStubServer()
		bundlerService = mock.StubBundlerService

		kitId := int64(9999)

		reader := bytes.NewReader([]byte("{}"))

		w := httptest.NewRecorder()
		req, err := http.NewRequest(http.MethodPost, fmt.Sprintf("/kits/%d/links", kitId), reader)

		assert.Nil(t, err)

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusNotFound, w.Code)
		assert.Equal(t, fmt.Sprintf("Kit %d not found", kitId), w.Body.String())
	})
}

func Test_RemoveKitLink(t *testing.T) {
	t.Run("should remove kit link", func(t *testing.T) {
		router := CreateStubServer()
		bundlerService = mock.StubBundlerService

		kit := mock.FakeKits[0]
		link := kit.Links[0]

		w := httptest.NewRecorder()
		uri := fmt.Sprintf("/kits/%d/links/%d", kit.ID, link.ID)
		req, err := http.NewRequest(http.MethodDelete, uri, nil)

		assert.Nil(t, err)

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusNoContent, w.Code)
	})
}

func Test_AddKitPart(t *testing.T) {
	t.Run("should add part to kit", func(t *testing.T) {
		router := CreateStubServer()
		bundlerService = mock.StubBundlerService

		kit := mock.FakeKits[0]
		part := mock.FakeParts[0]

		w := httptest.NewRecorder()
		uri := fmt.Sprintf("/kits/%d/parts/%d", kit.ID, part.ID)
		req, err := http.NewRequest(http.MethodPost, uri, nil)

		assert.Nil(t, err)

		router.ServeHTTP(w, req)

		var kitPart core.KitPart
		err = json.Unmarshal(w.Body.Bytes(), &kitPart)

		assert.Nil(t, err)
		assert.Equal(t, http.StatusOK, w.Code)
		assert.Equal(t, part.ID, kitPart.ID)
		assert.Equal(t, uint64(1), kitPart.Quantity)
	})

	t.Run("should use quantity value from query", func(t *testing.T) {
		router := CreateStubServer()
		bundlerService = mock.StubBundlerService

		kit := mock.FakeKits[0]
		part := mock.FakeParts[0]
		quantity := uint64(7)

		w := httptest.NewRecorder()
		uri := fmt.Sprintf("/kits/%d/parts/%d?quantity=%d", kit.ID, part.ID, quantity)
		req, err := http.NewRequest(http.MethodPost, uri, nil)

		assert.Nil(t, err)

		router.ServeHTTP(w, req)

		var kitPart core.KitPart
		err = json.Unmarshal(w.Body.Bytes(), &kitPart)

		assert.Nil(t, err)
		assert.Equal(t, http.StatusOK, w.Code)
		assert.Equal(t, part.ID, kitPart.ID)
		assert.Equal(t, quantity, kitPart.Quantity)
	})
}

func Test_RemoveKitPart(t *testing.T) {
	t.Run("should remove part from kit", func(t *testing.T) {
		router := CreateStubServer()
		bundlerService = mock.StubBundlerService

		kit := mock.FakeKits[0]
		part := mock.FakeParts[0]

		w := httptest.NewRecorder()
		uri := fmt.Sprintf("/kits/%d/parts/%d", kit.ID, part.ID)
		req, err := http.NewRequest(http.MethodDelete, uri, nil)

		assert.Nil(t, err)

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusNoContent, w.Code)
	})
}

func Test_UpdatePartQuantity(t *testing.T) {
	t.Run("should update part kit quantity", func(t *testing.T) {
		router := CreateStubServer()
		bundlerService = mock.StubBundlerService

		kit := mock.FakeKits[0]
		part := kit.Parts[0]
		newQty := part.Quantity * 2

		w := httptest.NewRecorder()
		uri := fmt.Sprintf("/kits/%d/parts/%d/%d", kit.ID, part.ID, newQty)
		req, err := http.NewRequest(http.MethodPut, uri, nil)

		assert.Nil(t, err)

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		kitPart := core.KitPart{}
		err = json.Unmarshal(w.Body.Bytes(), &kitPart)

		assert.Nil(t, err)
		assert.Equal(t, newQty, kitPart.Quantity)
	})
}
