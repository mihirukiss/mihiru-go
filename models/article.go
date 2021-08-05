package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type ArticleBaseFields struct {
	Author  string   `bson:"author" json:"author"`
	Title   string   `bson:"title" json:"title"`
	Summary string   `bson:"summary" json:"summary"`
	Tags    []string `bson:"tags" json:"tags"`
}

type ArticlePointFields struct {
	Ratting     int8  `bson:"ratting" json:"ratting"`
	PublishTime int64 `bson:"publishTime" json:"publishTime"`
	Hide        int8  `bson:"hide" json:"hide"`
}

type ArticleContentFields struct {
	Content string `bson:"content" json:"content"`
}

type ArticleAutoGenFields struct {
	ID      int64 `bson:"id" json:"id"`
	AddTime int64 `bson:"addTime" json:"addTime"`
	Version int32 `bson:"version" json:"version"`
}

type ArticleSearchParams struct {
	PageParams
	Keyword    string   `json:"keyword"`
	AllowTags  []string `json:"allowTags"`
	DenyTags   []string `json:"denyTags"`
	MaxRatting *int8    `json:"maxRatting"`
	ShowHide   *bool    `json:"showHide"`
}

type Article struct {
	ArticleAutoGenFields `bson:",inline"`
	ArticleBaseFields    `bson:",inline"`
	ArticleContentFields `bson:",inline"`
	ArticlePointFields   `bson:",inline"`
}

type ArticleWithObjectId struct {
	ObjectId primitive.ObjectID `bson:"_id"`
	Article  `bson:",inline"`
}

type ArticlePage struct {
	PageResult
	Data []*Article
}
