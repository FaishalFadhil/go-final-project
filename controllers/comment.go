package controllers

import (
	database "final-project/config/postgres"
	"final-project/helpers"
	"final-project/models"
	"net/http"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
)

func CreateComment(c *gin.Context) {
	db := database.GetDB()
	contentType := helpers.GetContentType(c)
	_, _ = db, contentType
	Comments := models.Comment{}
	userData := c.MustGet("userData").(jwt.MapClaims)
	userID := uint(userData["id"].(float64))

	if contentType == appJSON {
		c.ShouldBindJSON(&Comments)
	} else {
		c.ShouldBind(&Comments)
	}

	var photos models.Photo
	err := db.Preload("User").First(&photos, "id = ?", Comments.PhotoID).Error

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Bad Request",
			"message": err.Error(),
		})
		return
	}

	Comments.UserID = userID

	errCreate := db.Debug().Create(&Comments).Error

	if errCreate != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"err":     "Bad Request",
			"message": errCreate.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"data": Comments})
}

func GetAllComments(c *gin.Context) {
	db := database.GetDB()
	var Comments []models.Comment
	err := db.Preload("User").Preload("Photo").Find(&Comments).Error

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Bad Request",
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": Comments})
}

func GetOneComment(c *gin.Context) {
	db := database.GetDB()
	var Comments models.Comment
	err := db.Preload("User").Preload("Photo").First(&Comments, "id = ?", c.Param("commentId")).Error

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Bad Request",
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": Comments})
}

func UpdateComment(c *gin.Context) {
	db := database.GetDB()

	// Check data exist
	var Comments models.Comment

	err := db.Preload("User").Preload("Photo").First(&Comments, "id = ?", c.Param("commentId")).Error

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "request not found",
			"message": err.Error(),
		})
		return
	}

	var input models.Comment
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}

	errUpdate := db.Debug().Model(&Comments).Updates(input).Error

	if errUpdate != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Bad Request",
			"message": errUpdate.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": Comments,
	})
}

func DeleteComment(c *gin.Context) {
	db := database.GetDB()
	var Comments models.Comment

	err := db.Preload("User").Preload("Photo").First(&Comments, "id = ?", c.Param("commentId")).Error

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "request not found",
			"message": err.Error(),
		})
		return
	}

	errDelete := db.Debug().Delete(&Comments).Error

	if errDelete != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Bad Request",
			"message": errDelete.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Your comment has been successfully deleted",
	})
}
