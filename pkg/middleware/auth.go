package middleware

import (
	"scylla/pkg/config"
	"scylla/pkg/exception"
	"scylla/pkg/utils"
	"strings"

	"github.com/gin-gonic/gin"
)

func JwtMiddleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		conf := config.Get()
		var token string
		authorizationHeader := ctx.GetHeader("Authorization")
		fields := strings.Fields(authorizationHeader)

		if len(fields) == 2 && fields[0] == "Bearer" {
			token = fields[1]
		}

		if token == "" {
			panic(exception.NewUnauthorizedHandler("empty token"))
			return
		}

		if utils.IsTokenBlacklisted(token) {
			panic(exception.NewUnauthorizedHandler("token has been blacklisted"))
			return
		}

		sub, err := utils.ValidateToken(token, conf.Jwt.Secret)
		if err != nil {
			panic(exception.NewUnauthorizedHandler(err.Error()))
			//ctx.AbortWithStatusJSON(401, gin.H{"error": err.Error()})
			return
		}

		claims, ok := sub.(map[string]interface{})
		if !ok {
			panic(exception.NewUnauthorizedHandler("invalid token claims"))
			return
		}

		email, ok := claims["email"].(string)
		if !ok {
			panic(exception.NewUnauthorizedHandler("email not found in token claims"))
			return
		}

		ctx.Set("currentUser", email)
		ctx.Next()
	}
}
