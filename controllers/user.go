package controllers

import (
	// "fmt"
	database "final-project/config/postgres"
	"final-project/helpers"
	"final-project/models"
	"net/http"

	// "encoding/json"

	"github.com/gin-gonic/gin"
)

var (
	appJSON = "application/json"
)

func UserRegister(c *gin.Context) {
	db := database.GetDB()
	contentType := helpers.GetContentType(c)
	_, _ = db, contentType
	User := models.User{}

	if contentType == appJSON {
		c.ShouldBindJSON(&User)
	} else {
		c.ShouldBind(&User)
	}

	err := db.Debug().Create(&User).Error

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Bad Request",
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"id":       User.Id,
		"email":    User.Email,
		"username": User.Username,
		"age":      User.Age,
	})
}

func UserLogin(c *gin.Context) {
	db := database.GetDB()
	contentType := helpers.GetContentType(c)
	_, _ = db, contentType
	User := models.User{}
	password := ""

	if contentType == appJSON {
		c.ShouldBindJSON(&User)
	} else {
		c.ShouldBind(&User)
	}

	password = User.Password

	err := db.Debug().Where("email = ?", User.Email).Take(&User).Error

	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error":   "Unauthorized",
			"message": "invalid email/password",
		})
		return
	}

	comparePass := helpers.ComparePass([]byte(User.Password), []byte(password))

	if !comparePass {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error":   "Unauthorized",
			"message": "invalid email/password",
		})
		return
	}

	token := helpers.GenerateToken(User.Id, User.Email)

	c.JSON(http.StatusOK, gin.H{
		"token": token,
	})

}

func UpdateUser(c *gin.Context) {
	db := database.GetDB()

	// Check data exist
	var User models.User

	err := db.First(&User, "id = ?", c.Param("id")).Error

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "request not found",
			"message": err.Error(),
		})
		return
	}

	var input models.User
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}

	errUpdate := db.Debug().Model(&User).Updates(input).Error

	if errUpdate != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Bad Request",
			"message": errUpdate.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"id":       User.Id,
		"email":    User.Email,
		"username": User.Username,
		"age":      User.Age,
		"updateAt": User.UpdatedAt,
	})
}

func DeleteUser(c *gin.Context) {
	db := database.GetDB()
	var User models.User

	err := db.First(&User, "id = ?", c.Param("id")).Error

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "request not found",
			"message": err.Error(),
		})
		return
	}

	errDelete := db.Debug().Delete(&User).Error

	if errDelete != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Bad Request",
			"message": errDelete.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Your account has been successfully deleted",
	})
}
