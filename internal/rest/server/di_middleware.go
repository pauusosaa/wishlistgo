package server

import (
	"github.com/gin-gonic/gin"
	"github.com/nmarsollier/commongo/log"
	"github.com/nmarsollier/commongo/rst"
	"github.com/pauusosaa/wishlistgo/internal/di"
	"github.com/pauusosaa/wishlistgo/internal/env"
)

// DiInjectorMiddleware crea un Injector por request y lo guarda en el contexto
func DiInjectorMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		var deps di.Injector
		depParam, exists := c.Get("di")

		if !exists {
			logger := rst.GinLogger(c, env.Get().FluentURL, env.Get().ServerName)
			deps = di.NewInjector(logger)
			c.Set("di", deps)
		} else {
			deps = depParam.(di.Injector)
		}

		c.Next()

		if c.Request.Method != "OPTIONS" {
			deps.Logger().WithField(log.LOG_FIELD_HTTP_STATUS, c.Writer.Status()).Info("Completed")
		}
	}
}

// GinDi devuelve el Injector desde el contexto
func GinDi(c *gin.Context) di.Injector {
	return c.MustGet("di").(di.Injector)
}
