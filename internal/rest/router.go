package rest

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/pauusosaa/wishlistgo/internal/env"
	"github.com/pauusosaa/wishlistgo/internal/rest/server"
)

// Start levanta el servidor REST de Wishlist
func Start() {
	engine := server.Router()
	InitRoutes(engine)
	engine.Run(fmt.Sprintf(":%d", env.Get().Port))
}

// InitRoutes registra los endpoints REST
func InitRoutes(engine *gin.Engine) {
	initGetWishlist(engine)
	initPostWishlistArticle(engine)
	initDeleteWishlistArticle(engine)
	initPostWishlistArticleCart(engine)
}


