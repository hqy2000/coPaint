package route

import (
	"coPaint/config"
	"coPaint/model"
	"fmt"
	"github.com/gin-gonic/gin"
	"strconv"
)

func Default(c *gin.Context) {
	c.JSON(200, gin.H{
		"code": 200,
		"data": "Hello World!",
	})
}

func filter(vs []model.Painting, f func(model.Painting) bool) []model.Painting {
	vsf := make([]model.Painting, 0)
	for _, v := range vs {
		if f(v) {
			vsf = append(vsf, v)
		}
	}
	return vsf
}

func List(c *gin.Context) {
	ID, _ := strconv.Atoi(c.Query("paintingId"))
	fmt.Println(ID)

	c.JSON(200, gin.H{
		"code": 200,
		"data": filter(config.Paintings, func(painting model.Painting) bool {
			return painting.PaintingID == ID
		}),
	})
}

func Upload(c *gin.Context) {
	image := c.PostForm("image")
	fmt.Println(image)
	ID, _ := strconv.Atoi(c.PostForm("paintingId"))
	painting := model.Painting{
		Image: image,
		PaintingID: ID,
	}
	config.Paintings = append(config.Paintings, painting)
}