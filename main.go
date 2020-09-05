package main

import (
	"coPaint/route"
	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()
	r.GET("/", route.Default)
	r.GET("/paintings/list", route.List)
	r.Run()
}