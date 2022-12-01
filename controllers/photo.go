package controllers

import (
	database "final-project/config/postgres"
	"final-project/helpers"
	"final-project/models"
	"net/http"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
)

func CreatePhoto(c *gin.Context) {
	db := database.GetDB()
	contentType := helpers.GetContentType(c)
	_, _ = db, contentType
	Photos := models.Photo{}
	userData := c.MustGet("userData").(jwt.MapClaims)
	userID := uint(userData["id"].(float64))

	if contentType == appJSON {
		c.ShouldBindJSON(&Photos)
	} else {
		c.ShouldBind(&Photos)
	}

	Photos.UserID = userID

	err := db.Debug().Create(&Photos).Error

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"err":     "Bad Request",
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"data": Photos})
}

func GetAllPhotos(c *gin.Context) {
	db := database.GetDB()
	var photos []models.Photo
	err := db.Preload("User").Find(&photos).Error

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Bad Request",
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": photos})
}

func GetOnePhoto(c *gin.Context) {
	db := database.GetDB()
	var photos models.Photo
	err := db.Preload("User").First(&photos, "id = ?", c.Param("photoId")).Error

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Bad Request",
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": photos})
}

func UpdatePhoto(c *gin.Context) {
	db := database.GetDB()

	// Check data exist
	var Photos models.Photo

	err := db.Preload("User").First(&Photos, "id = ?", c.Param("photoId")).Error

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "request not found",
			"message": err.Error(),
		})
		return
	}

	var input models.Photo
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}

	errUpdate := db.Debug().Model(&Photos).Updates(input).Error

	if errUpdate != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Bad Request",
			"message": errUpdate.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": Photos,
	})
}

func DeletePhoto(c *gin.Context) {
	db := database.GetDB()
	var Photos models.Photo

	err := db.Preload("User").First(&Photos, "id = ?", c.Param("photoId")).Error

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "request not found",
			"message": err.Error(),
		})
		return
	}

	errDelete := db.Debug().Delete(&Photos).Error

	if errDelete != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Bad Request",
			"message": errDelete.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Your photo has been successfully deleted",
	})
}
