package controllers

import (
	"errors"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"log"
	"mihiru-go/models"
	"mihiru-go/services"
	"mihiru-go/util"
	"net/http"
	"strconv"
	"strings"
)

type VoiceController interface {
	AddVoice(c *gin.Context)
	UpdateVoice(c *gin.Context)
	DeleteVoice(c *gin.Context)
	LiverVoices(c *gin.Context)
}

type voiceController struct {
	service services.VoiceService
}

func NewVoiceController(service services.VoiceService) VoiceController {
	return voiceController{service: service}
}

func (v voiceController) AddVoice(c *gin.Context) {
	var voice models.VoiceBaseFields
	sortNoStr := c.PostForm("sort_no")
	if !checkEmptyParam(sortNoStr) {
		sortNo, err := strconv.ParseInt(sortNoStr, 10, 64)
		if err != nil {
			log.Println(err.Error())
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": "错误的排序号格式"})
			return
		}
		voice.SortNo = sortNo
	}
	voice.Liver = c.PostForm("liver")
	voice.Title = c.PostForm("title")
	voice.Category = c.PostForm("category")
	voice.Remark = c.PostForm("remark")
	if checkEmptyParam(voice.Liver) || checkEmptyParam(voice.Category) || checkEmptyParam(voice.Title) {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": "缺少必要参数"})
		return
	}

	file, err := c.FormFile("file")
	if err != nil {
		log.Println(err.Error())
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": "读取上传文件错误"})
		return
	}
	if file == nil || !strings.HasSuffix(strings.ToLower(file.Filename), ".mp3") {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": "请上传MP3格式音频文件"})
		return
	}
	err = v.service.AddVoice(&voice, file, c)
	if err != nil {
		util.ErrorResponse(c, err)
		return
	}
	c.Status(http.StatusOK)
}

func (v voiceController) UpdateVoice(c *gin.Context) {
	var voice models.VoiceBaseFields
	sortNoStr := c.PostForm("sort_no")
	if !checkEmptyParam(sortNoStr) {
		sortNo, err := strconv.ParseInt(sortNoStr, 10, 64)
		if err != nil {
			log.Println(err.Error())
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": "错误的排序号格式"})
			return
		}
		voice.SortNo = sortNo
	}
	voice.Liver = c.PostForm("liver")
	voice.Title = c.PostForm("title")
	voice.Category = c.PostForm("category")
	voice.Remark = c.PostForm("remark")
	if checkEmptyParam(voice.Liver) || checkEmptyParam(voice.Category) || checkEmptyParam(voice.Title) {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": "缺少必要参数"})
		return
	}

	id := c.Param("id")
	hex, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": "无效的ID格式"})
		return
	}
	file, err := c.FormFile("file")
	if err != nil && !errors.Is(err, http.ErrMissingFile) {
		log.Println(err.Error())
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": "读取上传文件错误"})
		return
	}
	if file != nil && !strings.HasSuffix(strings.ToLower(file.Filename), ".mp3") {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": "请上传MP3格式音频文件"})
		return
	}
	err = v.service.UpdateVoice(hex, &voice, file, c)
	if err != nil {
		util.ErrorResponse(c, err)
		return
	}
	c.Status(http.StatusOK)
}

func (v voiceController) DeleteVoice(c *gin.Context) {
	id := c.Param("id")
	hex, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": "无效的ID格式"})
		return
	}
	err = v.service.DeleteVoice(hex)
	if err != nil {
		util.ErrorResponse(c, err)
		return
	}
	c.Status(http.StatusOK)
}

func (v voiceController) LiverVoices(c *gin.Context) {
	liver := c.Param("liver")
	data, version, err := v.service.LiverVoices(liver)
	if err != nil {
		util.ErrorResponse(c, err)
		return
	}
	etag := strconv.FormatInt(version, 10)
	c.Header("ETag", etag)
	c.Header("Cache-Control", "public, max-age=0, must-revalidate")
	if match := c.GetHeader("If-None-Match"); match == etag {
		c.AbortWithStatus(http.StatusNotModified)
		return
	}
	c.JSON(http.StatusOK, data)
}

func checkEmptyParam(param string) bool {
	return len(strings.TrimSpace(param)) <= 0
}
