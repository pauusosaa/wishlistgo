package rest

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/pauusosaa/wishlistgo/internal/rest/server"
)

// initPostWishlistArticleCart registra POST /v1/wishlist/articles/:article_id/cart
func initPostWishlistArticleCart(engine *gin.Engine) {
	engine.POST(
		"/v1/wishlist/articles/:article_id/cart",
		moveWishlistArticleToCart,
	)
}

// moveWishlistArticleToCart mueve un artículo de la wishlist al carrito (por ahora solo lo elimina de la wishlist)
func moveWishlistArticleToCart(c *gin.Context) {
	userID := c.GetHeader("X-Demo-User-ID")
	if userID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "X-Demo-User-ID header requerido"})
		return
	}

	articleID := c.Param("article_id")

	deps := server.GinDi(c)
	wSvc := deps.WishlistService()

	err := wSvc.MoveToCart(userID, articleID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":               "Artículo agregado al carrito exitosamente (simulado)",
		"article_id":            articleID,
		"quantity":              1,
		"removed_from_wishlist": true,
	})
}


