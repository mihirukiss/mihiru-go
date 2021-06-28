package controllers

import (
	"github.com/gin-gonic/gin"
	"log"
	"mihiru-go/dto"
	"mihiru-go/models"
	"mihiru-go/services"
	"mihiru-go/util"
	"net/http"
	"strconv"
)

type ArticlesController interface {
	Add(c *gin.Context)
	Update(c *gin.Context)
	Get(c *gin.Context)
	Tags(c *gin.Context)
	Search(c *gin.Context)
}

type articlesController struct {
	articleService services.ArticleService
	userService    services.UserService
}

func NewArticlesController(articleService services.ArticleService, userService services.UserService) ArticlesController {
	return articlesController{articleService: articleService, userService: userService}
}

func (m articlesController) Add(c *gin.Context) {
	var articleDto dto.ArticleDto
	if err := c.BindJSON(&articleDto); err != nil {
		log.Println(err.Error())
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": "错误的参数格式"})
		return
	}
	articleVo, err := m.articleService.Add(&articleDto)
	if err != nil {
		util.ErrorResponse(c, err)
		return
	}
	c.JSON(http.StatusOK, articleVo)
}

func (m articlesController) Update(c *gin.Context) {
	var articleDto dto.ArticleDto
	if err := c.BindJSON(&articleDto); err != nil {
		return
	}
	id := c.Param("id")
	if id == "" {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": "缺少id参数"})
		return
	}
	intId, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		log.Println(err.Error())
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": "无效的id参数"})
		return
	}
	articleVo, err := m.articleService.Update(intId, &articleDto)
	if err != nil {
		util.ErrorResponse(c, err)
		return
	}
	c.JSON(http.StatusOK, articleVo)
}

func (m articlesController) Get(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": "缺少id参数"})
		return
	}
	intId, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		log.Println(err.Error())
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": "无效的id参数"})
		return
	}
	articleVo, err := m.articleService.Get(intId)
	if err != nil {
		util.ErrorResponse(c, err)
		return
	}
	if articleVo.Hide > 0 {
		if !m.checkIsAdmin(c) {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"message": "无权访问"})
			return
		}
	}
	etag := strconv.FormatInt(int64(articleVo.Version), 10)
	if c.Param("v") != "" {
		c.Header("Cache-Control", "public, max-age=31536000, must-revalidate")
	} else {
		c.Header("Cache-Control", "public, max-age=300, must-revalidate")
	}
	c.Header("ETag", etag)
	if match := c.GetHeader("If-None-Match"); match == etag {
		c.AbortWithStatus(http.StatusNotModified)
		return
	}
	c.JSON(http.StatusOK, articleVo)
}

func (m articlesController) Tags(c *gin.Context) {
	tags, err := m.articleService.Tags()
	if err != nil {
		util.ErrorResponse(c, err)
		return
	}
	c.JSON(http.StatusOK, tags)
}

func (m articlesController) Search(c *gin.Context) {
	var articleSearchParams models.ArticleSearchParams
	if err := c.BindJSON(&articleSearchParams); err != nil {
		return
	}
	if articleSearchParams.ShowHide != nil && *articleSearchParams.ShowHide && !m.checkIsAdmin(c) {
		c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"message": "无权访问"})
		return
	}
	result, err := m.articleService.Search(&articleSearchParams)
	if err != nil {
		util.ErrorResponse(c, err)
		return
	}
	c.JSON(http.StatusOK, result)
}

func (m articlesController) checkIsAdmin(c *gin.Context) bool {
	token := c.GetHeader("authorization")
	if token == "" {
		return false
	}
	user := m.userService.CheckToken(token)
	if user == nil {
		return false
	}
	hasPermission := false
	for _, userRole := range user.Roles {
		if userRole == "admin" {
			hasPermission = true
			break
		}
	}
	return hasPermission
}
