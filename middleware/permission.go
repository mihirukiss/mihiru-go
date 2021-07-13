package middleware

import (
	"github.com/gin-gonic/gin"
	"mihiru-go/services"
	"mihiru-go/vo"
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

func (p permissions) Login() gin.HandlerFunc {
	return func(c *gin.Context) {
		p.getLoginUser(c)
	}
}

func (p permissions) Roles(roles []string) gin.HandlerFunc {
	return func(c *gin.Context) {
		user := p.getLoginUser(c)
		if user == nil {
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

func (p permissions) getLoginUser(c *gin.Context) *vo.UserVo {
	token := c.GetHeader("authorization")
	if token == "" {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"message": "缺少验证信息"})
		return nil
	}
	user := p.userService.CheckToken(token)
	if user == nil {
		c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"message": "无效验证信息"})
		return nil
	}
	return user
}
