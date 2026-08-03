package requestutils

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

var ErrNoClaims = errors.New("token format invalid")

func GetUserIDFromRequest(c *gin.Context) (string, error) {
	claims, ok := c.Get("claims")
	if !ok {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"message": "Invalid token"})
		c.Abort()
		return "", ErrNoClaims
	}

	tokenInfo, ok := claims.(jwt.MapClaims)
	if !ok {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"message": "Invalid token"})
		c.Abort()
		return "", ErrNoClaims
	}

	uid, ok := tokenInfo["sub"].(string)
	if !ok {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"message": "Invalid token"})
		c.Abort()
		return "", ErrNoClaims
	}

	return uid, nil
}
