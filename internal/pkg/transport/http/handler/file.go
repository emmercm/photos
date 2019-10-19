package handler

import (
	"io/ioutil"
	"net/http"
	"os"

	"github.com/emmercm/photos/internal/pkg/model"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
)

// GetFiles responds with all Files
func GetFiles(db *gorm.DB) func(c *gin.Context) {
	return func(c *gin.Context) {
		var files []model.File

		if err := db.Find(&files).Error; err != nil {
			if !gorm.IsRecordNotFoundError(err) {
				c.JSON(http.StatusInternalServerError, gin.H{})
				return
			}
		}

		c.JSON(http.StatusOK, files)
	}
}

func getFile(db *gorm.DB, id string) (int, model.File) {
	var file model.File

	if err := db.Where("id = ?", id).First(&file).Error; err != nil {
		if gorm.IsRecordNotFoundError(err) {
			return http.StatusNotFound, file
		}
		return http.StatusInternalServerError, file
	}

	return http.StatusOK, file
}

// GetFile responds with a specific File
func GetFile(db *gorm.DB) func(c *gin.Context) {
	return func(c *gin.Context) {
		c.JSON(getFile(db, c.Param("id")))
	}
}

// GetFileImage responds with the image for a File
func GetFileImage(db *gorm.DB) func(c *gin.Context) {
	return func(c *gin.Context) {
		status, file := getFile(db, c.Param("id"))
		if status != http.StatusOK {
			c.JSON(status, gin.H{})
			return
		}

		f, err := os.Open(file.Path)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{})
			return
		}
		defer f.Close()

		b, err := ioutil.ReadAll(f)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{})
			return
		}

		c.Data(http.StatusOK, http.DetectContentType(b), b)
	}
}
