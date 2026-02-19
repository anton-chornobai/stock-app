package httphandler

import (
	tokenmanager "auth-service/internal/lib/jwt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

func ValidateJWT(secret []byte) gin.HandlerFunc {

	return func(c *gin.Context) {
		tokenStr, err := c.Cookie("token")

		if err != nil {
			c.JSON(http.StatusUnauthorized, map[string]any{
				"error": "missing or expired token",
			})
			//stops the middleware chain and doesnt go to the next handler
			c.Abort()
			return
		}

		claims := &tokenmanager.UserClaims{}

		token, err := jwt.ParseWithClaims(tokenStr, claims, func(t *jwt.Token) (any, error) {
			if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, jwt.ErrSignatureInvalid
			}

			return []byte(secret), nil
		})

		if err != nil || !token.Valid {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid token"})
			c.Abort()
			return
		}
		c.Set("user_id", claims.UserID)
		c.Set("role", claims.Role)
	}
}
