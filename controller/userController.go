package controller

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/deep/database"
	helper "github.com/deep/helpers"
	"github.com/deep/models"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

func HashPassword(password string) string {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	if err != nil {
		log.Panic(err)
	}
	return string(bytes)
}
func VerifyPassword(userPassword string, providedPassword string) (bool, string) {
	err := bcrypt.CompareHashAndPassword([]byte(userPassword), []byte(providedPassword))
	check := true
	msg := ""
	if err != nil {
		msg = fmt.Sprintf("email or password is incorrect")
		check = false
	}
	return check, msg
}
func Signup() gin.HandlerFunc {
	return func(c *gin.Context) {
		var user models.Users
		if err := c.BindJSON(&user); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		var users []models.Users
		db, err := database.ConnectionStablish()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
			return
		}
		password := HashPassword(user.Password)
		user.Password = password
		db.Find(&users, "email=?", user.Email)
		if len(users) > 0 {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "this email or phone already exist"})
			return
		}
		user.CreatedAt, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
		user.UpdatedAt, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
		token, refreshToken, err := helper.GenerateAllTokens(user.Email, user.First_Name, user.Last_Name, user.User_type, user.PhoneNumber)
		if len(users) > 0 {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "invalid request"})
			return
		}
		user.Token = token
		user.RefreshToken = refreshToken
		db.Create(&user)
	}
}
func Login() gin.HandlerFunc {
	return func(c *gin.Context) {
		var user models.Users
		var foundUser models.Users
		if err := c.BindJSON(&user); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		db, err := database.ConnectionStablish()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
			return
		}
		db.First(&foundUser, "email=?", user.Email)
		passwordIsValid, msg := VerifyPassword(foundUser.Password, user.Password)
		if passwordIsValid != true {
			c.JSON(http.StatusInternalServerError, gin.H{"error": msg})
			return
		}
		token, refershToken, err := helper.GenerateAllTokens(foundUser.Email, foundUser.First_Name, foundUser.Last_Name, foundUser.User_type, foundUser.PhoneNumber)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		foundUser.Token = token
		foundUser.RefreshToken = refershToken
		c.JSON(http.StatusOK, foundUser)
	}
}

func GetUser() gin.HandlerFunc {
	return func(c *gin.Context) {
		userId := c.Param("user_id")
		db, err := database.ConnectionStablish()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
			return
		}
		var user models.Users
		db.Find(&user, "email=?", userId)
		c.JSON(http.StatusOK, user)
	}
}
func UpdateUser() gin.HandlerFunc {
	return func(c *gin.Context) {
		userId := c.Param("user_id")
		var user models.Users
		var foundUser models.Users
		if err := c.BindJSON(&user); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		db, err := database.ConnectionStablish()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
			return
		}
		db.First(&foundUser, "email=?", userId)
		user.UpdatedAt, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
		db.Model(foundUser).Updates(models.Users{PhoneNumber: user.PhoneNumber, UpdatedAt: user.UpdatedAt})
		db.First(&foundUser, "email=?", userId)
		c.JSON(http.StatusOK, foundUser)
	}
}
func DeleteUser() gin.HandlerFunc {
	return func(c *gin.Context) {
		userId := c.Param("user_id")
		email, ok := c.Get("email")
		if !ok {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "invalid request"})
			return
		}
		if email != userId {
			c.JSON(http.StatusBadRequest, gin.H{"error": "you not authorised to delete this account"})
		}
		db, err := database.ConnectionStablish()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
			return
		}
		db.Where("email", email).Delete(&models.Users{})
		c.JSON(http.StatusOK, gin.H{"msg": "successfully deleted"})
	}
}
