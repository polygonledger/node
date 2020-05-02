package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()

	r.LoadHTMLGlob("templates/*")

	r.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"data": "hello world"})
	})

	//curl --data "name=bla" http://localhost:8080/name
	r.GET("/name", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.tmpl", []string{"a", "b", "c"})
	})
	r.POST("/name", func(c *gin.Context) {
		//name := c.Params.ByName("name")
		//c.JSON(http.StatusOK, gin.H{"name": name})
		c.JSON(http.StatusOK, gin.H{"name": "name"})
	})

	r.Run()
}
