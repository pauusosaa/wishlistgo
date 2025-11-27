package rest

import (
	"github.com/gin-gonic/gin"
	"github.com/nmarsollier/commongo/rst"
	"github.com/nmarsollier/commongo/security"
	"github.com/pauusosaa/wishlistgo/internal/rest/server"
)

func initGetWishlist(engine *gin.Engine) {
	engine.GET(
		"/v1/wishlist",
		server.ValidateAuthentication,
		getWishlist,
	)
}

func getWishlist(c *gin.Context) {
	user := c.MustGet("user").(security.User)
	tokenString := c.MustGet("tokenString").(string)

	deps := server.GinDi(c)
	wSvc := deps.WishlistAppService()

	w, err := wSvc.GetWishlist(user.ID, tokenString)
	if err != nil {
		rst.AbortWithError(c, err)
		return
	}

	c.JSON(200, w)
}


