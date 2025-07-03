package handlers

import (
	"crypto/rand"
	"encoding/base64"
	"net/http"
	"pos/config"
	"pos/database"
	"pos/models"
	"strconv"
	"time"

	"pos/utils"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
	"golang.org/x/crypto/bcrypt"
)

type AuthClaims struct {
	jwt.RegisteredClaims
	Issuer string `json:"issuer"`
}

type CreateUserInput struct {
	Username string  `json:"username" binding:"required"`
	Name     string  `json:"name" binding:"required"`
	Password *string `json:"password" binding:"required"`
}

type LoginInput struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type LoginResponse struct {
	Token        string `json:"token"`
	RefreshToken string `json:"refresh_token"`
}

type RefreshTokenInput struct {
	RefreshToken string `json:"refresh_token"`
}

type ErrorResponse struct {
	Message string `json:"message"`
}

func generateRefreshToken() (string, error) {
	b := make([]byte, 32)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return base64.URLEncoding.EncodeToString(b), nil
}

// @Summary Register a new user
// @Description Register a new user with username, name, and password.
// @Tags Auth
// @Accept  json
// @Produce  json
// @Param   user     body    CreateUserInput     true        "User registration info"
// @Success 200 {object} models.User
// @Failure 400 {object} ErrorResponse
// @Failure 406 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /auth/register [post]
func Register(c *gin.Context) {
	var data CreateUserInput

	if err := c.ShouldBindJSON(&data); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Message: err.Error(),
		})
		return
	}
	isDup, err := utils.IsDuplicate[models.User](database.DB, "username", data.Username, nil)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Message: "Database error",
		})
		return
	}
	if isDup {
		c.JSON(406, ErrorResponse{
			Message: "Username already exists",
		})
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

// @Summary Login a user
// @Description Login a user with username and password to get access and refresh tokens.
// @Tags Auth
// @Accept  json
// @Produce  json
// @Param   credentials body LoginInput true "User login credentials"
// @Success 200 {object} LoginResponse
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Router /auth/login [post]
func Login(c *gin.Context) {
	var data LoginInput

	if err := c.ShouldBindJSON(&data); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Message: err.Error(),
		})
		return
	}

	var user models.User

	database.DB.Where("username = ?", data.Username).First(&user)

	if user.ID == 0 {
		c.JSON(http.StatusNotFound, ErrorResponse{
			Message: "User not found",
		})
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(data.Password)); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Message: "Incorrect password",
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
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Message: "Could not login",
		})
		return
	}

	refreshToken, err := generateRefreshToken()
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Message: "Could not generate refresh token",
		})
		return
	}

	user.RefreshToken = refreshToken
	database.DB.Save(&user)

	c.JSON(http.StatusOK, LoginResponse{
		Token:        tokenString,
		RefreshToken: refreshToken,
	})
}

// @Summary Refresh access token
// @Description Refresh an expired access token using a refresh token.
// @Tags Auth
// @Accept  json
// @Produce  json
// @Param   token    body    RefreshTokenInput true "Refresh token"
// @Success 200 {object} LoginResponse
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Router /auth/refresh [post]
func RefreshToken(c *gin.Context) {
	var data struct {
		RefreshToken string `json:"refresh_token"`
	}

	if err := c.BindJSON(&data); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Message: "Invalid request body"})
		return
	}

	var user models.User
	database.DB.Where("refresh_token = ?", data.RefreshToken).First(&user)

	if user.ID == 0 {
		c.JSON(http.StatusUnauthorized, ErrorResponse{Message: "Invalid refresh token"})
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
		c.JSON(http.StatusInternalServerError, ErrorResponse{Message: "Could not generate access token"})
		return
	}

	// Generate new refresh token
	newRefreshToken, err := generateRefreshToken()
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Message: "Could not generate new refresh token"})
		return
	}

	user.RefreshToken = newRefreshToken
	database.DB.Save(&user)
	c.JSON(http.StatusOK, LoginResponse{
		Token:        tokenString,
		RefreshToken: newRefreshToken,
	})
}
