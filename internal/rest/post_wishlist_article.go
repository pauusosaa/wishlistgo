package rest

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/nmarsollier/commongo/errs"
	"github.com/nmarsollier/commongo/rst"
	"github.com/nmarsollier/commongo/security"
	"github.com/pauusosaa/wishlistgo/internal/rest/server"
)

type addToWishlistBody struct {
	ArticleID string  `json:"article_id" binding:"required"`
	Notes     *string `json:"notes"`
}

func initPostWishlistArticle(engine *gin.Engine) {
	engine.POST(
		"/v1/wishlist/article",
		server.ValidateAuthentication,
		addToWishlist,
	)
}

func addToWishlist(c *gin.Context) {
	user := c.MustGet("user").(security.User)
	tokenString := c.MustGet("tokenString").(string)

	var body addToWishlistBody
	if err := c.ShouldBindJSON(&body); err != nil {
		rst.AbortWithError(c, err)
		return
	}

	deps := server.GinDi(c)
	wSvc := deps.WishlistAppService()

	item, err := wSvc.AddToWishlist(user.ID, body.ArticleID, body.Notes, tokenString)
	if err != nil {
		if validationErr, ok := err.(*errs.ValidationErr); ok {
			for _, msg := range validationErr.Messages {
				if msg.Path == "article_id" {
					c.JSON(http.StatusConflict, gin.H{"error": msg.Message})
					return
				}
			}
		}
		rst.AbortWithError(c, err)
		return
	}

	c.JSON(201, item)
}


