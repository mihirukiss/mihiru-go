package vo

import "mihiru-go/models"

type ArticleListVo struct {
	models.ArticleAutoGenFields
	models.ArticleBaseFields
	models.ArticlePointFields
}

type ArticlePageVo struct {
	models.PageResult
	Data *[]ArticleListVo `json:"data"`
}

type ArticleVo struct {
	ArticleListVo
	models.ArticleContentFields
}
