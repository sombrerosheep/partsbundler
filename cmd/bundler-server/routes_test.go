package main

import (
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
		router.ServeHTTP(w, req)

		assert.Nil(t, err)

		var actualParts []core.Part

		err = json.Unmarshal(w.Body.Bytes(), &actualParts)

		assert.Nil(t, err)
		assert.Equal(t, http.StatusOK, w.Code)
		assert.Equal(t, mock.FakeParts[:], actualParts)
	})
}

func Test_GetPart(t *testing.T) {
	t.Run("should each part", func(t *testing.T) {
		router := CreateStubServer()
		bundlerService = mock.StubBundlerService

		for _, v := range mock.FakeParts {

			w := httptest.NewRecorder()
			req, err := http.NewRequest(http.MethodGet, fmt.Sprintf("/parts/%d", v.ID), nil)
			router.ServeHTTP(w, req)

			assert.Nil(t, err)

			var part core.Part

			err = json.Unmarshal(w.Body.Bytes(), &part)

			assert.Nil(t, err)
			assert.Equal(t, http.StatusOK, w.Code)
			assert.Equal(t, v, part)
		}
	})
}

func Test_GetAllKits(t *testing.T) {
	t.Run("should return kits", func(t *testing.T) {
		router := CreateStubServer()
		bundlerService = mock.StubBundlerService

		w := httptest.NewRecorder()
		req, err := http.NewRequest(http.MethodGet, "/kits", nil)
		router.ServeHTTP(w, req)

		assert.Nil(t, err)

		var actualKits []core.Kit

		err = json.Unmarshal(w.Body.Bytes(), &actualKits)

		assert.Nil(t, err)
		assert.Equal(t, http.StatusOK, w.Code)
		assert.Equal(t, mock.FakeKits[:], actualKits)
	})
}
