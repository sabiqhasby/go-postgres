package auth

import (
	"learn2/models"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	jwt "github.com/golang-jwt/jwt/v5"
)

const (
	USER     = "admin"
	PASSWORD = "Password123"
	SECRET   = "secret"
)

func LoginHandler(c *gin.Context) {
	var user models.Credential

	err := c.Bind(&user)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "bad request"})
		return
	}

	if user.Username != USER {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "user invalid"})
		return
	} else if user.Password != PASSWORD {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "user invalid"})
		return
	}

	//token
	claim := &jwt.RegisteredClaims{
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Minute * 1)),
		Issuer:    "test",
		IssuedAt:  jwt.NewNumericDate(time.Now()),
	}

	sign := jwt.NewWithClaims(jwt.SigningMethodHS256, claim)
	token, err := sign.SignedString([]byte(SECRET))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		c.Abort()
	}
	c.JSON(http.StatusOK, gin.H{"message": "success", "token": token})

}
