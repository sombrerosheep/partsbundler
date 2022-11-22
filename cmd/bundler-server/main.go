package main

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/sombrerosheep/partsbundler/internal/sqlite"
	"github.com/sombrerosheep/partsbundler/pkg/service"
)

const bundlerDBPath string = "../../data/partsbundler.db"

var bundlerService *service.BundlerService = nil

func InitBundlerService(dbPath string) error {
	svc, err := sqlite.CreateSqliteService(dbPath)
	if err != nil {
		return err
	}

	bundlerService = svc

	return nil
}

func GetBundlerService() *service.BundlerService {
	return bundlerService
}

func RegisterEndpoints(router *gin.Engine, endpoints []Endpoint) {
	for _, v := range endpoints {
		switch v.method {
		case http.MethodGet:
			router.GET(v.path, v.handler)
		case http.MethodPost:
			router.POST(v.path, v.handler)
		case http.MethodDelete:
			router.DELETE(v.path, v.handler)
		case http.MethodPut:
			router.PUT(v.path, v.handler)
		default:
			fmt.Printf("Unsupported method '%s' for endpoint %#v", v.method, v)
		}
	}
}

func main() {
	fmt.Println("Hello")

	err := InitBundlerService(bundlerDBPath)
	if err != nil {
		fmt.Printf("Error iniializing service: %s\n", err)
		return
	}

	router := gin.Default()
	RegisterEndpoints(router, endpoints)

	err = router.Run(":3000")
	if err != nil {
		fmt.Printf("Server exited with error: %s", err)
	}
}
