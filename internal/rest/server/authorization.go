package server

import (
	"github.com/gin-gonic/gin"
	"github.com/nmarsollier/commongo/errs"
	"github.com/nmarsollier/commongo/log"
	"github.com/nmarsollier/commongo/rst"
	"github.com/nmarsollier/commongo/security"
)

// ValidateAuthentication valida el token JWT y guarda el usuario en el contexto
func ValidateAuthentication(c *gin.Context) {
	user, err := validateToken(c)
	if err != nil {
		c.Error(err)
		c.Abort()
		return
	}

	deps := GinDi(c)
	deps.Logger().WithField(log.LOG_FIELD_USER_ID, user.ID)
}

func validateToken(c *gin.Context) (*security.User, error) {
	tokenString, err := rst.GetHeaderToken(c)
	if err != nil {
		return nil, errs.Unauthorized
	}
	c.Set("tokenString", tokenString)

	deps := GinDi(c)
	user, err := deps.SecurityService().Validate(tokenString)
	c.Set("user", *user)
	if err != nil {
		return nil, errs.Unauthorized
	}

	return user, nil
}

