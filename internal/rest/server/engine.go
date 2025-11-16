package server

import (
	"github.com/gin-gonic/gin"
	"github.com/nmarsollier/wishlistgo/internal/di"
)

var engine *gin.Engine = nil

// Router inicializa el engine de Gin con el middleware de DI
func Router() *gin.Engine {
	if engine != nil {
		return engine
	}

	engine = gin.Default()
	engine.Use(DiInjectorMiddleware())

	return engine
}

// GinDi devuelve el Injector desde el contexto
func GinDi(c *gin.Context) di.Injector {
	return c.MustGet("di").(di.Injector)
}


