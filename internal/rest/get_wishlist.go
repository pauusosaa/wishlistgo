package rest

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/pauusosaa/wishlistgo/internal/rest/server"
)

// initGetWishlist registra GET /v1/wishlist
func initGetWishlist(engine *gin.Engine) {
	engine.GET(
		"/v1/wishlist",
		getWishlist,
	)
}

// getWishlist obtiene la wishlist del usuario
func getWishlist(c *gin.Context) {
	// En una versión real, el user_id vendría del token JWT (Auth Service)
	userID := c.GetHeader("X-Demo-User-ID")
	if userID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "X-Demo-User-ID header requerido"})
		return
	}

	deps := server.GinDi(c)
	wSvc := deps.WishlistService()

	w, err := wSvc.GetWishlist(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, w)
}


