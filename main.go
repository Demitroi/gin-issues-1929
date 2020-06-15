package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"

	"github.com/gin-gonic/gin"
)

func main() {
	// Start router configuration
	var router *gin.Engine
	r, _ := os.LookupEnv("GINROUTER")
	switch r {
	case "nolog":
		log.Println("nolog")
		gin.DefaultWriter = ioutil.Discard
		gin.DefaultErrorWriter = ioutil.Discard
		router = gin.Default()
	case "go-api":
		log.Println("go-api")
		router = gin.New()
		router.Use(Recover())
	default:
		log.Println("default")
		router = gin.Default()
	}
	// Route with panic
	router.GET("/panic", func(c *gin.Context) {
		panic("FINISH!")
	})
	// Route ok
	router.GET("/ok", func(c *gin.Context) {
		c.JSON(200, gin.H{"msg": "OK"})
	})
	// Select the port
	var addr = ":8080"
	if port, ok := os.LookupEnv("GINPORT"); ok {
		addr = fmt.Sprintf(":%s", port)
	}
	// Start router
	log.Fatal(router.Run(addr))
}
