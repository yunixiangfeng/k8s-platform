package middle

import (
	"k8s-platform/utils"
	"net/http"

	"github.com/gin-gonic/gin"
)

func JWTAuth() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		// 对登录接口放行
		if len(ctx.Request.URL.String()) >= 10 && ctx.Request.URL.String()[0:10] == "/api/login" {
			ctx.Next()
		} else {
			// 处理验证逻辑
			token := ctx.Request.Header.Get("Authorization")
			if token == "" {
				ctx.JSON(http.StatusBadRequest, gin.H{
					"msg":  "请求未携带token，无权限访问",
					"data": nil,
				})
				ctx.Abort()
				return
			}
			// 解析token内容
			claims, err := utils.JWTToken.ParseToken(token)
			if err != nil {
				// token过期错误
				if err.Error() == "TokenExpired" {
					ctx.JSON(http.StatusBadRequest, gin.H{
						"msg":  "授权已过期",
						"data": nil,
					})
					ctx.Abort()
					return
				}
				// 其他解析错误
				ctx.JSON(http.StatusBadRequest, gin.H{
					"msg":  err.Error(),
					"data": nil,
				})
				ctx.Abort()
				return
			}
			ctx.Set("claims", claims)
			ctx.Next()
		}

	}
}
