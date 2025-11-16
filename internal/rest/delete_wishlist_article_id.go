package rest

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/pauusosaa/wishlistgo/internal/rest/server"
)

// initDeleteWishlistArticle registra DELETE /v1/wishlist/article/:article_id
func initDeleteWishlistArticle(engine *gin.Engine) {
	engine.DELETE(
		"/v1/wishlist/article/:article_id",
		deleteWishlistArticle,
	)
}

// deleteWishlistArticle elimina un artículo de la wishlist
func deleteWishlistArticle(c *gin.Context) {
	userID := c.GetHeader("X-Demo-User-ID")
	if userID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "X-Demo-User-ID header requerido"})
		return
	}

	articleID := c.Param("article_id")

	deps := server.GinDi(c)
	wSvc := deps.WishlistService()

	found, err := wSvc.RemoveFromWishlist(userID, articleID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if !found {
		c.JSON(http.StatusNotFound, gin.H{"error": "artículo no está en la wishlist"})
		return
	}

	c.Status(http.StatusNoContent)
}


