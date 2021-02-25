package main

import (
	"encoding/json"
	"fmt"

	"github.com/gin-gonic/gin"
)

type flow struct {
	Url      string          `json:"url"`
	Name     string          `json:"name"`
	Function gin.HandlerFunc `json:"-"`
}

func main() {

	flows := []flow{
		flow{
			Url:      "/lights_on",
			Name:     "Lights On",
			Function: lightson,
		},
		flow{
			Url:      "/lights_off",
			Name:     "Lights Off",
			Function: lightsoff,
		},
	}

	router := gin.Default()

	router.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})

	router.GET("/flows", func(c *gin.Context) {
		json.NewEncoder(c.Writer).Encode(&flows)
	})

	for _, f := range flows {
		router.GET(f.Url, f.Function)
	}

	router.Run(":6666")
}

func lightson(c *gin.Context) {
	fmt.Println("lightson")
}

func lightsoff(c *gin.Context) {
	fmt.Println("lightsoff")
}
