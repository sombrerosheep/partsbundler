package main

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/sombrerosheep/partsbundler/internal/sqlite"
)

const dbPath string = "../../data/partsbundler.db"

func main() {
	fmt.Println("Hello")

	svc, err := sqlite.CreateSqliteService(dbPath)
	if err != nil {
		fmt.Printf("unable to create service: %s", err)
		return
	}

	router := gin.Default()

	router.GET("/parts", func(c *gin.Context) {
		parts, err := svc.Parts.GetAll()
		if err != nil {
			c.String(http.StatusInternalServerError, err.Error())
		}

		c.JSON(http.StatusOK, parts)
	})

	router.GET("/kits", func(c *gin.Context) {
		kits, err := svc.Kits.GetAll()
		if err != nil {
			c.String(http.StatusInternalServerError, err.Error())
		}

		c.JSON(http.StatusOK, kits)
	})

	err = router.Run(":3000")
	if err != nil {
		fmt.Printf("Server exited with error: %s", err)
	}
}
