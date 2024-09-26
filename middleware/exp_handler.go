package middleware

import (
	"gin-exp-middleware/exp"
	"github.com/gin-gonic/gin"
	"net/http"
)

func ExpMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()

		for _, err := range c.Errors {
			if e := exp.IsHttpExp(err.Err); e != nil {
				c.AbortWithStatusJSON(e.Status, e)
			} else if e := exp.IsValidationExp(err.Err); e != nil {
				c.AbortWithStatusJSON(e.Status, e)
			} else {
				c.AbortWithStatus(http.StatusInternalServerError)
			}
		}
	}
}
