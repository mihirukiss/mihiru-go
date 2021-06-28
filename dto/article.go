package dto

import "mihiru-go/models"

type ArticleDto struct {
	models.ArticleBaseFields
	models.ArticleContentFields
	Hide        *int8  `json:"hide"`
	Ratting     *int8  `json:"ratting"`
	PublishTime *int64 `json:"publishTime"`
	AutoFormat  *bool  `json:"autoFormat"`
	Indent      *bool  `json:"indent"`
}
