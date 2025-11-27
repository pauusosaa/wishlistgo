package server

import (
	"github.com/gin-gonic/gin"
	"github.com/nmarsollier/commongo/rst"
)

var engine *gin.Engine = nil

// Router inicializa el engine de Gin con el middleware de DI
func Router() *gin.Engine {
	if engine != nil {
		return engine
	}

	engine = gin.Default()
	engine.Use(DiInjectorMiddleware())
	engine.Use(rst.ErrorHandler)

	return engine
}
