package route

import (
	"coPaint/config"
	"github.com/gin-gonic/gin"

)

func Default(c *gin.Context) {
	c.JSON(200, gin.H{
		"code": 200,
		"data": "Hello World!",
	})
}

func List(c *gin.Context) {
	c.JSON(200, gin.H{
		"code": 200,
		"data": config.Paintings,
	})
}

func Login(c *gin.Context) {

}