package middleware

import (
	"github.com/gin-gonic/gin"
	"mihiru-go/services"
	"net/http"
)

type Permissions interface {
	Roles(roles []string) gin.HandlerFunc
}

type permissions struct {
	userService services.UserService
}

func NewPermissions(userService services.UserService) Permissions {
	return permissions{userService: userService}
}

func (p permissions) Roles(roles []string) gin.HandlerFunc {
	return func(c *gin.Context) {
		token := c.GetHeader("authorization")
		if token == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"message": "缺少验证信息"})
			return
		}
		user := p.userService.CheckToken(token)
		if user == nil {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"message": "无效验证信息"})
			return
		}
		for _, userRole := range user.Roles {
			for _, role := range roles {
				if userRole == role {
					c.Next()
					return
				}
			}
		}
		c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"message": "无权访问"})
	}
}
