package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/medlcago/gin-exp-middleware/exp"
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
