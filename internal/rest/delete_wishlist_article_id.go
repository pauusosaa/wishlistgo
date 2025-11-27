package rest

import (
	"github.com/gin-gonic/gin"
	"github.com/nmarsollier/commongo/errs"
	"github.com/nmarsollier/commongo/rst"
	"github.com/nmarsollier/commongo/security"
	"github.com/pauusosaa/wishlistgo/internal/rest/server"
)

func initDeleteWishlistArticle(engine *gin.Engine) {
	engine.DELETE(
		"/v1/wishlist/article/:article_id",
		server.ValidateAuthentication,
		deleteWishlistArticle,
	)
}

func deleteWishlistArticle(c *gin.Context) {
	user := c.MustGet("user").(security.User)
	articleID := c.Param("article_id")

	deps := server.GinDi(c)
	wSvc := deps.WishlistAppService()

	found, err := wSvc.RemoveFromWishlist(user.ID, articleID)
	if err != nil {
		rst.AbortWithError(c, err)
		return
	}

	if !found {
		rst.AbortWithError(c, errs.NotFound)
		return
	}

	c.Status(204)
}


