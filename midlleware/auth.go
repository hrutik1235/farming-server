package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/hrutik1235/farming-server/utils"
)

func GateValidateUser() gin.HandlerFunc {
	return func(c *gin.Context) {
		userId := c.GetHeader("user_id")

		if userId == "" {
			c.JSON(http.StatusUnauthorized, utils.NewHttpError(c, "Unauthorized", http.StatusUnauthorized))
			c.Abort() // Stop further processing
			return
		}

		// If userId is present, continue to the next handler
		c.Next()
	}
}
