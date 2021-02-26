package main

import (
	"encoding/json"
	"net"

	"github.com/gin-gonic/gin"
)

const lights_on_command = "{\"id\": 1, \"method\": \"set_power\", \"params\": [\"on\", \"smooth\", 300]}\r\n"
const lights_off_command = "{\"id\": 1, \"method\": \"set_power\", \"params\": [\"off\", \"smooth\", 300]}\r\n"

func lights_ips() [4]string {
	return [4]string{
		"192.168.0.220",
		"192.168.0.221",
		"192.168.0.222",
		"192.168.0.223",
	}
}

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

	router.GET("/flows", func(c *gin.Context) {
		json.NewEncoder(c.Writer).Encode(&flows)
	})

	for _, f := range flows {
		router.GET(f.Url, f.Function)
	}

	router.Run(":6666")
}

func lightson(c *gin.Context) {
	all_lights_command(lights_on_command)
}

func lightsoff(c *gin.Context) {
	all_lights_command(lights_off_command)
}

func all_lights_command(command string) {
	for _, ip := range lights_ips() {
		go send_command_to_light(ip, command)
	}
}

func send_command_to_light(ip string, command string) {
	tcpAddr, _ := net.ResolveTCPAddr("tcp", ip+":55443")
	conn, _ := net.DialTCP("tcp", nil, tcpAddr)
	conn.Write([]byte(command))
	conn.CloseWrite()
}
