package controllers

import (
	database "final-project/config/postgres"
	"final-project/helpers"
	"final-project/models"
	"net/http"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
)

func CreateSocialMedia(c *gin.Context) {
	db := database.GetDB()
	contentType := helpers.GetContentType(c)
	_, _ = db, contentType
	SocialMedias := models.SocialMedia{}
	userData := c.MustGet("userData").(jwt.MapClaims)
	userID := uint(userData["id"].(float64))

	if contentType == appJSON {
		c.ShouldBindJSON(&SocialMedias)
	} else {
		c.ShouldBind(&SocialMedias)
	}

	SocialMedias.UserID = userID

	err := db.Debug().Create(&SocialMedias).Error

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"err":     "Bad Request",
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"data": SocialMedias})
}

func GetAllSocialMedias(c *gin.Context) {
	db := database.GetDB()
	var SocialMedias []models.SocialMedia
	err := db.Preload("User").Find(&SocialMedias).Error

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Bad Request",
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": SocialMedias})
}

func GetOneSocialMedia(c *gin.Context) {
	db := database.GetDB()
	var SocialMedias models.SocialMedia
	err := db.Preload("User").First(&SocialMedias, "id = ?", c.Param("socialMediaId")).Error

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Bad Request",
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": SocialMedias})
}

func UpdateSocialMedia(c *gin.Context) {
	db := database.GetDB()

	// Check data exist
	var SocialMedias models.SocialMedia

	err := db.Preload("User").First(&SocialMedias, "id = ?", c.Param("socialMediaId")).Error

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "request not found",
			"message": err.Error(),
		})
		return
	}

	var input models.SocialMedia
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}

	errUpdate := db.Debug().Model(&SocialMedias).Updates(input).Error

	if errUpdate != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Bad Request",
			"message": errUpdate.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": SocialMedias,
	})
}

func DeleteSocialMedia(c *gin.Context) {
	db := database.GetDB()
	var SocialMedias models.SocialMedia

	err := db.Preload("User").First(&SocialMedias, "id = ?", c.Param("socialMediaId")).Error

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "request not found",
			"message": err.Error(),
		})
		return
	}

	errDelete := db.Debug().Delete(&SocialMedias).Error

	if errDelete != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Bad Request",
			"message": errDelete.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Your social media account has been successfully deleted",
	})
}
