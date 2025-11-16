package server

import (
	"github.com/gin-gonic/gin"
	"github.com/nmarsollier/wishlistgo/internal/di"
)

// DiInjectorMiddleware crea un Injector por request y lo guarda en el contexto
func DiInjectorMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// En esta primera versión no usamos logger ni tracing,
		// simplemente creamos el inyector por request.
		deps := di.NewInjector()
		c.Set("di", deps)

		c.Next()
	}
}


