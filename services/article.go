package services

import (
	"bytes"
	"mihiru-go/database"
	"mihiru-go/dto"
	"mihiru-go/models"
	"mihiru-go/util"
	"mihiru-go/vo"
	"net/http"
	"strings"
	"time"
	"unicode/utf8"
)

type ArticleService interface {
	Add(articleDto *dto.ArticleDto) (*vo.ArticleVo, error)
	Update(id int64, articleDto *dto.ArticleDto) (*vo.ArticleVo, error)
	Get(id int64) (*vo.ArticleVo, error)
	Search(articleSearchParams *models.ArticleSearchParams) (*vo.ArticlePageVo, error)
	Tags() ([]string, error)
}

type articleService struct {
	db        database.ArticleDatabase
	tagsCache []string
}

func NewArticleService(db database.ArticleDatabase) ArticleService {
	return articleService{db: db}
}

func (a articleService) Add(articleDto *dto.ArticleDto) (*vo.ArticleVo, error) {
	article := new(models.Article)
	article.Author = articleDto.Author
	article.Title = articleDto.Title
	if articleDto.Hide != nil && *articleDto.Hide > 0 {
		article.Hide = 1
	} else {
		article.Hide = 0
	}
	if articleDto.Ratting != nil {
		article.Ratting = *articleDto.Ratting
	} else {
		article.Ratting = 0
	}
	article.AddTime = time.Now().UnixNano() / 1e6
	if articleDto.PublishTime != nil {
		article.PublishTime = *articleDto.PublishTime
	} else {
		article.PublishTime = article.AddTime
	}
	article.Version = 1
	tagsLen := len(articleDto.Tags)
	formattedTags := make([]string, tagsLen)
	if tagsLen > 0 {
		for i := range articleDto.Tags {
			formattedTags[i] = strings.Title(articleDto.Tags[i])
		}
		article.Tags = formattedTags
	}
	handleArticleContent(articleDto, article)
	err := a.db.InsertArticle(article)
	if err != nil {
		util.LogError(err)
		return nil, vo.NewErrorWithHttpStatus("添加文章失败, 请稍后重试", http.StatusInternalServerError)
	}
	if tagsLen > 0 {
		a.tagsCache = a.tagsCache[:0]
	}
	return convertToArticleVo(article), nil
}

func (a articleService) Update(id int64, articleDto *dto.ArticleDto) (*vo.ArticleVo, error) {
	article, err := a.db.GetArticle(id)
	if err != nil {
		util.LogError(err)
		return nil, vo.NewErrorWithHttpStatus("查询数据失败, 请稍后重试", http.StatusInternalServerError)
	}
	if article == nil {
		return nil, vo.NewErrorWithHttpStatus("无效的文章ID", http.StatusNotFound)
	}
	if articleDto.Author != "" {
		article.Author = articleDto.Author
	}
	if articleDto.Title != "" {
		article.Title = articleDto.Title
	} else {
		articleDto.Title = article.Title
	}
	if articleDto.Hide != nil {
		if *articleDto.Hide > 0 {
			article.Hide = 1
		} else {
			article.Hide = 0
		}
	}
	if articleDto.Ratting != nil {
		article.Ratting = *articleDto.Ratting
	}
	if articleDto.PublishTime != nil {
		article.PublishTime = *articleDto.PublishTime
	}
	article.Version++
	tagsLen := len(articleDto.Tags)
	formattedTags := make([]string, tagsLen)
	if tagsLen > 0 {
		for i := range articleDto.Tags {
			formattedTags[i] = strings.Title(articleDto.Tags[i])
		}
		article.Tags = formattedTags
	}
	if articleDto.Content != "" {
		handleArticleContent(articleDto, &article.Article)
	}
	if articleDto.Summary != "" {
		article.Summary = articleDto.Summary
	}
	err = a.db.UpdateArticle(article)
	if err != nil {
		util.LogError(err)
		return nil, vo.NewErrorWithHttpStatus("更新数据失败, 请稍后重试", http.StatusInternalServerError)
	}
	if tagsLen > 0 {
		a.tagsCache = a.tagsCache[:0]
	}
	return convertToArticleVo(&article.Article), nil
}

func (a articleService) Get(id int64) (*vo.ArticleVo, error) {
	article, err := a.db.GetArticle(id)
	if err != nil {
		util.LogError(err)
		return nil, vo.NewErrorWithHttpStatus("查询数据失败, 请稍后重试", http.StatusInternalServerError)
	}
	if article == nil {
		return nil, vo.NewErrorWithHttpStatus("无效的文章ID", http.StatusNotFound)
	}
	return convertToArticleVo(&article.Article), nil
}

func (a articleService) Search(articleSearchParams *models.ArticleSearchParams) (*vo.ArticlePageVo, error) {
	articles, err := a.db.SearchArticle(articleSearchParams)
	if err != nil {
		util.LogError(err)
		return nil, vo.NewErrorWithHttpStatus("查询数据失败, 请稍后重试", http.StatusInternalServerError)
	}
	pageVo := new(vo.ArticlePageVo)
	pageVo.PageResult = articles.PageResult
	var data []vo.ArticleListVo
	for i := range articles.Data {
		data = append(data, *convertToArticleListVo(articles.Data[i]))
	}
	pageVo.Data = &data
	return pageVo, nil
}

func (a articleService) Tags() ([]string, error) {
	if len(a.tagsCache) == 0 {
		allTag, err := a.db.ListAllTag()
		if err != nil {
			util.LogError(err)
			return nil, vo.NewErrorWithHttpStatus("查询数据失败, 请稍后重试", http.StatusInternalServerError)
		}
		a.tagsCache = allTag
	}
	return a.tagsCache, nil
}

func handleArticleContent(articleDto *dto.ArticleDto, article *models.Article) {
	if articleDto.AutoFormat != nil && *articleDto.AutoFormat {
		var buffer bytes.Buffer
		var summaryBuffer bytes.Buffer
		summaryLength := 0
		summaryBufferLength := 0
		if utf8.RuneCountInString(articleDto.Title) > 13 {
			summaryLength = 19 * 4
		} else {
			summaryLength = 19 * 5
		}
		indent := articleDto.Indent != nil && *articleDto.Indent
		trimContent := strings.TrimSpace(articleDto.Content)
		splitContent := strings.Split(trimContent, "\n")
		emptyLineCount := 0
		summaryFin := false

		buffer.WriteString("<p")
		if indent {
			buffer.WriteString(" class=\"indent\"")
		}
		buffer.WriteString(">")
		for i := range splitContent {
			trimLine := strings.TrimSpace(splitContent[i])
			if trimLine == "" {
				emptyLineCount++
			} else {
				if i == 0 {
					//Do Nothing
				} else if emptyLineCount > 1 {
					buffer.WriteString("</p><br/><p")
					if indent {
						buffer.WriteString(" class=\"indent\"")
					}
					buffer.WriteString(">")
				} else if emptyLineCount > 0 {
					buffer.WriteString("</p><p")
					if indent {
						buffer.WriteString(" class=\"indent\"")
					}
					buffer.WriteString(">")
				} else {
					buffer.WriteString("<br/>")
				}
				buffer.WriteString(trimLine)
				if !summaryFin {
					if i > 0 {
						summaryBuffer.WriteString(" ")
					}
					summaryBuffer.WriteString(trimLine)
					summaryBufferLength += utf8.RuneCountInString(trimLine)
					if summaryBufferLength > summaryLength {
						summaryFin = true
					}
				}
				emptyLineCount = 0
			}
		}
		buffer.WriteString("</p>")
		article.Content = buffer.String()
		article.Summary = summaryBuffer.String()
		if summaryFin {
			article.Summary = strings.TrimSpace(string(([]rune(article.Summary))[:summaryLength-3])) + "..."
		}
	} else {
		article.Content = articleDto.Content
		if articleDto.Summary != "" {
			article.Summary = articleDto.Summary
		}
	}
}

func convertToArticleVo(article *models.Article) *vo.ArticleVo {
	articleVo := new(vo.ArticleVo)
	articleVo.ArticleContentFields = article.ArticleContentFields
	articleVo.ArticleAutoGenFields = article.ArticleAutoGenFields
	articleVo.ArticleBaseFields = article.ArticleBaseFields
	articleVo.ArticlePointFields = article.ArticlePointFields
	return articleVo
}

func convertToArticleListVo(article *models.Article) *vo.ArticleListVo {
	articleListVo := new(vo.ArticleListVo)
	articleListVo.ArticleAutoGenFields = article.ArticleAutoGenFields
	articleListVo.ArticleBaseFields = article.ArticleBaseFields
	articleListVo.ArticlePointFields = article.ArticlePointFields
	return articleListVo
}
