package rest

import (
	"github.com/gin-gonic/gin"
	"github.com/nmarsollier/commongo/rst"
	"github.com/nmarsollier/commongo/security"
	"github.com/pauusosaa/wishlistgo/internal/rest/server"
)

func initPostWishlistArticleCart(engine *gin.Engine) {
	engine.POST(
		"/v1/wishlist/articles/:article_id/cart",
		server.ValidateAuthentication,
		moveWishlistArticleToCart,
	)
}

func moveWishlistArticleToCart(c *gin.Context) {
	user := c.MustGet("user").(security.User)
	tokenString := c.MustGet("tokenString").(string)
	articleID := c.Param("article_id")

	deps := server.GinDi(c)
	wSvc := deps.WishlistAppService()

	err := wSvc.MoveToCart(user.ID, articleID, tokenString)
	if err != nil {
		rst.AbortWithError(c, err)
		return
	}

	c.JSON(200, gin.H{
		"message":               "Artículo agregado al carrito exitosamente",
		"article_id":            articleID,
		"quantity":              1,
		"removed_from_wishlist": true,
	})
}


