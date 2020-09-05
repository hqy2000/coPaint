package route

import (
	"coPaint/config"
	"coPaint/model"
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

func Upload(c *gin.Context) {
	image := c.PostForm("image")
	painting := model.Painting{
		ID: 2,
		Name: "lang.java.NullPointerException",
		Image: image,
	}
	config.Paintings = append(config.Paintings, painting)
}