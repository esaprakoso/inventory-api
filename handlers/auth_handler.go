package handlers

import (
	"crypto/rand"
	"encoding/base64"
	"inventory/config"
	"inventory/database"
	"inventory/models"
	"net/http"
	"strconv"
	"time"

	"inventory/utils"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
	"golang.org/x/crypto/bcrypt"
)

type AuthClaims struct {
	jwt.RegisteredClaims
	Issuer string `json:"issuer"`
}

func generateRefreshToken() (string, error) {
	b := make([]byte, 32)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return base64.URLEncoding.EncodeToString(b), nil
}

func Register(c *gin.Context) {
	type CreateUserInput struct {
		Username string  `json:"username" binding:"required"`
		Name     string  `json:"name" binding:"required"`
		Password *string `json:"password" binding:"required"`
	}
	var data CreateUserInput

	if err := c.ShouldBindJSON(&data); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}
	isDup, err := utils.IsDuplicate[models.User](database.DB, "username", data.Username, nil)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Database error"})
		return
	}
	if isDup {
		c.JSON(406, gin.H{"message": "Username already exists"})
		return
	}

	password, _ := bcrypt.GenerateFromPassword([]byte(*data.Password), 14)

	user := models.User{
		Username: data.Username,
		Name:     data.Name,
		Password: string(password),
		Role:     models.RoleUser, // Set default role
	}

	database.DB.Create(&user)

	c.JSON(http.StatusOK, user)
}

func Login(c *gin.Context) {
	type LoginInput struct {
		Username string `json:"username" binding:"required"`
		Password string `json:"password" binding:"required"`
	}
	var data LoginInput

	if err := c.ShouldBindJSON(&data); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	var user models.User

	database.DB.Where("username = ?", data.Username).First(&user)

	if user.ID == 0 {
		c.JSON(http.StatusNotFound, gin.H{
			"message": "User not found",
		})
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(data.Password)); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Incorrect password",
		})
		return
	}

	claims := AuthClaims{
		Issuer: strconv.Itoa(int(user.ID)),
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * 24)),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenString, err := token.SignedString([]byte(config.LoadConfig("JWT_SECRET")))

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Could not login",
		})
		return
	}

	refreshToken, err := generateRefreshToken()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Could not generate refresh token"})
		return
	}

	user.RefreshToken = refreshToken
	database.DB.Save(&user)

	c.JSON(http.StatusOK, gin.H{
		"token":         tokenString,
		"refresh_token": refreshToken,
	})
}

func RefreshToken(c *gin.Context) {
	var data struct {
		RefreshToken string `json:"refresh_token"`
	}

	if err := c.BindJSON(&data); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid request body"})
		return
	}

	var user models.User
	database.DB.Where("refresh_token = ?", data.RefreshToken).First(&user)

	if user.ID == 0 {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "Invalid refresh token"})
		return
	}

	// Generate new access token
	claims := AuthClaims{
		Issuer: strconv.Itoa(int(user.ID)),
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * 24)),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(config.LoadConfig("JWT_SECRET")))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Could not generate access token"})
		return
	}

	// Generate new refresh token
	newRefreshToken, err := generateRefreshToken()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Could not generate new refresh token"})
		return
	}

	user.RefreshToken = newRefreshToken
	database.DB.Save(&user)

	c.JSON(http.StatusOK, gin.H{
		"token":         tokenString,
		"refresh_token": newRefreshToken,
	})
}
