// https://github.com/daheige/go-api/blob/e41d1cd23bf4d1061eb91a44d1071b2ddcb44472/app/middleware/log.go
package main

import (
	"strings"
	"os"
	"net"
	"github.com/gin-gonic/gin"
)


// Recover 请求处理中遇到异常或panic捕获
func Recover() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		defer func() {
			if err := recover(); err != nil {

				// Check for a broken connection, as it is not really a
				// condition that warrants a panic stack trace.
				var brokenPipe bool
				if ne, ok := err.(*net.OpError); ok {
					if se, ok := ne.Err.(*os.SyscallError); ok {

						errMsg := strings.ToLower(se.Error())


						if strings.Contains(errMsg, "broken pipe") ||
							strings.Contains(errMsg, "reset by peer") ||
							strings.Contains(errMsg, "request headers: small read buffer") ||
							strings.Contains(errMsg, "unexpected EOF") ||
							strings.Contains(errMsg, "i/o timeout") {
							brokenPipe = true
						}
					}
				}

				// 是否是 brokenPipe类型的错误
				// 如果是该类型的错误，就不需要返回任何数据给客户端
				// 代码参考gin recovery.go RecoveryWithWriter方法实现
				// If the connection is dead, we can't write a status to it.
				// if broken pipe,return nothing.
				if brokenPipe {
					// ctx.Error(err.(error)) // nolint: errcheck
					ctx.Abort()
					return
				}

				//响应状态
				ctx.AbortWithStatusJSON(500, gin.H{
					"code":    500,
					"message": "server error",
				})
			}
		}()

		ctx.Next()
	}

}