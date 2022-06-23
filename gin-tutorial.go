package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type album struct {
	ID     string  `json:"id"`
	Title  string  `json:"title"`
	Artist string  `json:"artist"`
	Price  float32 `json:"price"`
}

var albums = []album{
	{"1", "Strong Wave", "James", 5.32},
	{"2", "Destiny", "Blacks", 4.99},
	{"3", "Love", "Annie", 10.01},
}

func test() {
	routers := gin.Default()
	routers.GET("/albums", getAlbums)
	routers.POST("/albums", postAlbums)
	routers.GET("/albums/:id", getAlbumById)

	routers.Run(":8000")

}

func getAlbums(c *gin.Context) {
	c.IndentedJSON(http.StatusOK, albums)
}

func postAlbums(c *gin.Context) {
	var newAlbum album

	err := c.BindJSON(&newAlbum)

	if err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "request json wrong param"})
	}

	albums = append(albums, newAlbum)

	c.IndentedJSON(http.StatusCreated, newAlbum)
}

func getAlbumById(c *gin.Context) {
	id := c.Param("id")

	for _, a := range albums {
		if a.ID == id {
			c.IndentedJSON(http.StatusOK, a)
			return
		}
	}

	c.IndentedJSON(http.StatusNotFound, gin.H{"message": "Album not found"})
}
