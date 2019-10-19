package handler

import (
	"net/http"

	"github.com/emmercm/photos/internal/pkg/model"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
)

// GetAlbums responds with all Albums
func GetAlbums(db *gorm.DB) func(c *gin.Context) {
	return func(c *gin.Context) {
		var albums []model.Album

		if err := db.Find(&albums).Error; err != nil {
			if !gorm.IsRecordNotFoundError(err) {
				c.JSON(http.StatusInternalServerError, gin.H{})
				return
			}
		}

		c.JSON(http.StatusOK, albums)
	}
}

// GetAlbum responds with a specific Album
func GetAlbum(db *gorm.DB) func(c *gin.Context) {
	return func(c *gin.Context) {
		var album model.Album
		id := c.Param("id")

		if err := db.Where("id = ?", id).Find(&album).Error; err != nil {
			if gorm.IsRecordNotFoundError(err) {
				c.JSON(http.StatusNotFound, gin.H{})
				return
			}
			c.JSON(http.StatusInternalServerError, gin.H{})
			return
		}

		c.JSON(http.StatusOK, album)
	}
}

// GetAlbumFiles responds with all Files for an Album
func GetAlbumFiles(db *gorm.DB) func(c *gin.Context) {
	return func(c *gin.Context) {
		id := c.Param("id")

		var album model.Album
		if err := db.Where("id = ?", id).Preload("Files").Find(&album).Error; err != nil {
			if gorm.IsRecordNotFoundError(err) {
				c.JSON(http.StatusNotFound, gin.H{})
				return
			}
			c.JSON(http.StatusInternalServerError, gin.H{})
			return
		}

		c.JSON(http.StatusOK, album.Files)
	}
}
