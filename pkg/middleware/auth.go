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
		var token string
		authorizationHeader := ctx.GetHeader("Authorization")
		fields := strings.Fields(authorizationHeader)

		if len(fields) == 2 && fields[0] == "Bearer" {
			token = fields[1]
		}

		if token == "" {
			panic(exception.NewUnauthorizedError("empty token"))
			return
		}

		if utils.IsTokenBlacklisted(token) {
			panic(exception.NewUnauthorizedError("token has been blacklisted"))
			return
		}

		config, err := config.LoadConfig(".")
		if err != nil {
			panic(exception.NewInternalServerError(err.Error()))
			return
		}

		sub, err := utils.ValidateToken(token, config.TokenSecret)
		if err != nil {
			panic(exception.NewUnauthorizedError(err.Error()))
			//ctx.AbortWithStatusJSON(401, gin.H{"error": err.Error()})
			return
		}

		claims, ok := sub.(map[string]interface{})
		if !ok {
			panic(exception.NewUnauthorizedError("invalid token claims"))
			return
		}

		email, ok := claims["email"].(string)
		if !ok {
			panic(exception.NewUnauthorizedError("email not found in token claims"))
			return
		}

		ctx.Set("currentUser", email)
		ctx.Next()
	}
}
