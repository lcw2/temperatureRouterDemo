package main

import (
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
)

func main() {
	var temperature string
	var humidity string
	r := gin.Default()
	r.GET("/", func(c *gin.Context) {
		c.String(http.StatusOK, "service is running")
	})

	r.POST("/environment", func(c *gin.Context) {
		data := make(map[string]interface{})
		err := c.BindJSON(&data)
		if err != nil {
			log.Println(err)
			c.JSON(200, gin.H{
				"errcode": 400,
			})
			return
		}
		temperature = data["temperature"].(string)
		humidity = data["humidity"].(string)
		c.JSON(http.StatusOK, gin.H{
			"status": "ok",
		})
	})

	r.GET("/environment", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"temperature": temperature,
			"humidity":    humidity,
		})
	})
	_ = r.Run()
}
