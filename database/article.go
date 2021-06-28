package database

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"mihiru-go/models"
	"mihiru-go/util"
	"strings"
)

const collectionNameArticle = "article"

type ArticleDatabase interface {
	InsertArticle(article *models.Article) error
	UpdateArticle(article *models.ArticleWithObjectId) error
	GetArticle(id int64) (*models.ArticleWithObjectId, error)
	SearchArticle(articleSearchParams *models.ArticleSearchParams) (*models.ArticlePage, error)
	ListAllTag() ([]string, error)
}

func (d *MongoDatabase) SearchArticle(articleSearchParams *models.ArticleSearchParams) (*models.ArticlePage, error) {
	var pageSize int64
	var pageIndex int64
	if articleSearchParams.PageSize != nil {
		pageSize = *articleSearchParams.PageSize
	} else {
		pageSize = 10
	}
	if articleSearchParams.PageIndex != nil {
		pageIndex = *articleSearchParams.PageIndex
	} else {
		pageIndex = 0
	}
	skip := pageSize * pageIndex
	filter := bson.D{}
	if articleSearchParams.ShowHide == nil || !*articleSearchParams.ShowHide {
		filter = append(filter, bson.E{Key: "hide", Value: int8(0)})
	}
	if articleSearchParams.MaxRatting != nil {
		filter = append(filter, bson.E{Key: "ratting", Value: bson.D{
			{"$lte", *articleSearchParams.MaxRatting},
		}})
	}
	keyword := strings.TrimSpace(articleSearchParams.Keyword)
	if keyword != "" {
		var builder strings.Builder
		runeKeyword := []rune(keyword)
		for i := 0; i < len(runeKeyword); i++ {
			char := string(runeKeyword[i])
			if util.TextInArray(char, d.escapeStrings) {
				builder.WriteString("\\")
			}
			builder.WriteString(char)
		}
		regex := primitive.Regex{Pattern: builder.String(), Options: "i"}
		filter = append(filter, bson.E{Key: "$or", Value: bson.A{
			bson.D{{"title", regex}},
			bson.D{{"content", regex}},
			bson.D{{"author", regex}},
		}})
	}
	if len(articleSearchParams.AllowTags) > 0 {
		filter = append(filter, bson.E{Key: "tags", Value: bson.M{"$in": articleSearchParams.AllowTags}})
	}
	if len(articleSearchParams.DenyTags) > 0 {
		filter = append(filter, bson.E{Key: "tags", Value: bson.M{"$nin": articleSearchParams.DenyTags}})
	}
	collection := d.DB.Collection(collectionNameArticle)
	count, err := collection.CountDocuments(context.Background(), filter)
	if err != nil {
		return nil, err
	}
	cursor, err := collection.
		Find(context.Background(), filter,
			&options.FindOptions{
				Skip:  &skip,
				Sort:  bson.D{bson.E{Key: "id", Value: -1}},
				Limit: &pageSize,
			})
	if err != nil {
		return nil, err
	}
	defer CloseCursor(cursor, context.Background())

	var data []*models.Article
	for cursor.Next(context.Background()) {
		var article *models.Article
		if err := cursor.Decode(&article); err != nil {
			return nil, err
		}
		data = append(data, article)
	}
	articlePage := new(models.ArticlePage)
	articlePage.Data = data
	articlePage.PageSize = &pageSize
	articlePage.PageIndex = &pageIndex
	articlePage.Count = count
	articlePage.PageCount = count / pageSize
	if count%pageSize > 0 {
		articlePage.PageCount++
	}
	return articlePage, nil
}

func (d *MongoDatabase) InsertArticle(article *models.Article) error {
	collection := d.DB.Collection(collectionNameArticle)
	var maxIdArticle *models.Article
	err := collection.FindOne(context.Background(), bson.D{}, &options.FindOneOptions{
		Sort:       bson.D{bson.E{Key: "id", Value: -1}},
		Projection: bson.M{"_id": 0, "id": 1},
	}).Decode(&maxIdArticle)
	var id int64
	if err != nil && err != mongo.ErrNoDocuments {
		return err
	} else if err == nil {
		id = maxIdArticle.ID
	}
	article.ID = id + 1
	if _, err = collection.InsertOne(context.Background(), article); err != nil {
		return err
	}
	return nil
}

func (d *MongoDatabase) UpdateArticle(article *models.ArticleWithObjectId) error {
	collection := d.DB.Collection(collectionNameArticle)
	_, err := collection.UpdateByID(context.Background(), article.ObjectId, article)
	if err != nil {
		return err
	}
	return nil
}

func (d *MongoDatabase) GetArticle(id int64) (*models.ArticleWithObjectId, error) {
	var article *models.ArticleWithObjectId
	err := d.DB.Collection(collectionNameArticle).
		FindOne(context.Background(), bson.D{{Key: "id", Value: id}}).
		Decode(&article)
	if err != nil && err != mongo.ErrNoDocuments {
		return nil, err
	}
	return article, nil
}

func (d *MongoDatabase) ListAllTag() ([]string, error) {
	var tags []string
	values, err := d.DB.Collection(collectionNameArticle).Distinct(context.Background(), "tags", bson.D{})
	if err != nil {
		return nil, err
	}
	for _, value := range values {
		tags = append(tags, value.(string))
	}
	return tags, nil
}
