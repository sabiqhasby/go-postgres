package middleware

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

const (
	SECRET = "secret"
)

func AuthValid(c *gin.Context) {
	tokenString := c.Request.Header.Get("Authorization")

	if tokenString == "" {
		c.JSON(http.StatusBadRequest, gin.H{"message": "token required"})
		c.Abort()
		return
	}

	token, err := jwt.Parse(tokenString, func(t *jwt.Token) (interface{}, error) {
		if _, invalid := t.Method.(*jwt.SigningMethodHMAC); !invalid {
			return nil, fmt.Errorf("invalid token ", t.Header["alg"])
		}

		return []byte(SECRET), nil
	})

	// if token != nil && err == nil {
	// 	fmt.Println("Token verified")
	// 	c.Next()
	// } else {
	// 	c.JSON(http.StatusUnauthorized, gin.H{"message": "not authorized", "error": err.Error()})
	// 	c.Abort()
	// }

	if token == nil || err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "not authorized", "error": "Token verification failed"})
		c.Abort()
		return
	} else {
		fmt.Println("Token verified")
		c.Next()
	}

}
