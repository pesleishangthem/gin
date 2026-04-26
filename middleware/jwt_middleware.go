package middleware

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
)

func JWTMiddleware(v *Validator) gin.HandlerFunc {

	return func(c *gin.Context) {
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}
		if strings.Contains(c.Request.URL.Path, "/health") {
			c.JSON(http.StatusOK, gin.H{"status": "ok"})
			c.AbortWithStatus(200)
			return
		}
		authHeader := c.GetHeader("Authorization")

		if authHeader == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "missing token"})
			return
		}

		tokenString := strings.TrimPrefix(authHeader, "Bearer ")

		token, err := v.Validate(tokenString)

		if err != nil || !token.Valid {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "invalid token"})
			return
		}

		claims := token.Claims.(jwt.MapClaims)
		var user User
		ubyte, err := json.Marshal(claims)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		err = json.Unmarshal(ubyte, &user)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		roles := make([]string, 0)
		for _, v := range user.Roles {
			if strings.Contains(v, "default-roles") || strings.Contains(v, "uma_authorization") {
				continue
			}
			roles = append(roles, v)
		}
		user.Roles = roles
		c.Set("user", user)
		c.Set("token", authHeader)

		c.Next()
	}
}
