package middlewares

import (
	"net/http"
	"strings"
	"task-backend/internal/res"
	my_utils "task-backend/utils"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")

		if authHeader == "" || authHeader == "null" || authHeader == "undefined" {
			userID := uuid.New().String()
			newToken, err := my_utils.GenerateJWT(userID)
			if err != nil {
				c.AbortWithStatusJSON(http.StatusInternalServerError, res.ErrorResponse{
					Message: "Failed to generate token",
					Status:  http.StatusInternalServerError,
					Error:   err.Error(),
				})
				return
			}

			c.Header("Authorization", "Bearer "+newToken)
			c.Set("userID", userID)
			c.Next()
			return
		}

		tokenParts := strings.Split(authHeader, " ")
		if len(tokenParts) != 2 || tokenParts[0] != "Bearer" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, res.ErrorResponse{
				Message: "Unauthorized",
				Status:  http.StatusUnauthorized,
				Error:   "invalid token format",
			})
			return
		}

		tokenString := tokenParts[1]
		claims, err := my_utils.ValidateJWT(tokenString)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, res.ErrorResponse{
				Message: "Unauthorized",
				Status:  http.StatusUnauthorized,
				Error:   err.Error(),
			})
			return
		}
		// extract user_id from JWT claims
		if userID, ok := claims["user_id"].(string); ok {
			c.Set("userID", userID)
		} else {
			c.AbortWithStatusJSON(http.StatusUnauthorized, res.ErrorResponse{
				Message: "Invalid token claims",
				Status:  http.StatusUnauthorized,
				Error:   "user_id claim missing or invalid",
			})
			return
		}

		c.Next()
	}
}
