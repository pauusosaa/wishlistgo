package rest

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/nmarsollier/wishlistgo/internal/rest/server"
)

type addToWishlistBody struct {
	ArticleID string  `json:"article_id" binding:"required"`
	Notes     *string `json:"notes"`
}

// initPostWishlistArticle registra POST /v1/wishlist/article
func initPostWishlistArticle(engine *gin.Engine) {
	engine.POST(
		"/v1/wishlist/article",
		addToWishlist,
	)
}

// addToWishlist agrega un artículo a la wishlist del usuario
func addToWishlist(c *gin.Context) {
	userID := c.GetHeader("X-Demo-User-ID")
	if userID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "X-Demo-User-ID header requerido"})
		return
	}

	var body addToWishlistBody
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	deps := server.GinDi(c)
	wSvc := deps.WishlistService()

	item, err := wSvc.AddToWishlist(userID, body.ArticleID, body.Notes)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, item)
}


