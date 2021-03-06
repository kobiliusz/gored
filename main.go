package main

import (
	"encoding/json"
	"net"
	"sync"

	"github.com/gin-gonic/gin"
	"github.com/kobiliusz/go-miio"
)

const lights_on_command = "{\"id\": 1, \"method\": \"set_power\", \"params\": [\"on\", \"smooth\", 300, 1 ]}\r\n"
const lights_off_command = "{\"id\": 1, \"method\": \"set_power\", \"params\": [\"off\", \"smooth\", 300]}\r\n"
const full_brightness_command = "{\"id\": 2, \"method\": \"set_bright\", \"params\": [100, \"smooth\", 300]}\r\n"
const colortemp_command = "{\"method\":\"props\",\"params\":{\"ct\":6500}}\r\n"
const rgbmode_command = "{\"id\": 3, \"method\": \"set_power\", \"params\": [\"on\", \"smooth\", 300, 2]}\r\n"
const colorblue_command = "{\"id\": 4, \"method\": \"start_cf\", \"params\": [1, 1, \"300, 1, 21247, 100\"]}\r\n"

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
		flow{
			Url:      "/lights_blue",
			Name:     "Lights Blue",
			Function: lightsblue,
		},
		flow{
			Url:      "/air_fav",
			Name:     "Air Full",
			Function: airheart,
		},
		flow{
			Url:      "/air_auto",
			Name:     "Air Auto",
			Function: airauto,
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
	all_lights_command(full_brightness_command)
	all_lights_command(colortemp_command)
}

func lightsoff(c *gin.Context) {
	all_lights_command(lights_off_command)
}

func lightsblue(c *gin.Context) {
	all_lights_command(rgbmode_command)
	all_lights_command(full_brightness_command)
	all_lights_command(colorblue_command)
}

func all_lights_command(command string) {
	var wg sync.WaitGroup
	wg.Add(len(lights_ips()))
	for _, ip := range lights_ips() {
		go func(ip string, command string) {
			defer wg.Done()
			tcpAddr, _ := net.ResolveTCPAddr("tcp", ip+":55443")
			conn, _ := net.DialTCP("tcp", nil, tcpAddr)
			conn.Write([]byte(command))
			conn.Close()
		}(ip, command)
	}
	wg.Wait()
}

func airheart(c *gin.Context) {
	val := "favorite"
	go air_command("set_mode", []interface{}{val})
}

func airauto(c *gin.Context) {
	val := "auto"
	go air_command("set_mode", []interface{}{val})
}

func air_command(command string, data []interface{}) {
	var air miio.XiaomiDevice
	air.Start(air_ip, air_token, miio.DefaultPort)
	air.SendCommand(command, data, false, 1)
	air.Stop()

}
