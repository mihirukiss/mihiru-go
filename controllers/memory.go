package controllers

import (
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"log"
	"mihiru-go/models"
	"mihiru-go/services"
	"mihiru-go/util"
	"net/http"
	"strconv"
)

type MemoryController interface {
	AddDynamic(c *gin.Context)
	UpdateDynamic(c *gin.Context)
	AddLive(c *gin.Context)
	UpdateLive(c *gin.Context)
	Days(c *gin.Context)
	Day(c *gin.Context)
}

type memoryController struct {
	service services.MemoryService
}

func NewMemoryController(service services.MemoryService) MemoryController {
	return memoryController{service: service}
}

func (m memoryController) AddDynamic(c *gin.Context) {
	var dynamic models.Dynamic
	if err := c.BindJSON(&dynamic); err != nil {
		log.Println(err.Error())
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": "错误的参数格式"})
		return
	}
	dynamicVo, err := m.service.AddDynamic(&dynamic)
	if err != nil {
		util.ErrorResponse(c, err)
		return
	}
	c.JSON(http.StatusOK, dynamicVo)
}

func (m memoryController) UpdateDynamic(c *gin.Context) {
	var dynamic models.Dynamic
	if err := c.BindJSON(&dynamic); err != nil {
		log.Println(err.Error())
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": "错误的参数格式"})
		return
	}
	id := c.Param("id")
	hex, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": "无效的ID格式"})
		return
	}
	dynamicVo, err := m.service.UpdateDynamic(hex, &dynamic)
	if err != nil {
		util.ErrorResponse(c, err)
		return
	}
	c.JSON(http.StatusOK, dynamicVo)
}

func (m memoryController) AddLive(c *gin.Context) {
	var live models.Live
	if err := c.BindJSON(&live); err != nil {
		log.Println(err.Error())
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": "错误的参数格式"})
		return
	}
	liveVo, err := m.service.AddLive(&live)
	if err != nil {
		util.ErrorResponse(c, err)
		return
	}
	c.JSON(http.StatusOK, liveVo)
}

func (m memoryController) UpdateLive(c *gin.Context) {
	var live models.Live
	if err := c.BindJSON(&live); err != nil {
		log.Println(err.Error())
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": "错误的参数格式"})
		return
	}
	id := c.Param("id")
	hex, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": "无效的ID格式"})
		return
	}
	liveVo, err := m.service.UpdateLive(hex, &live)
	if err != nil {
		util.ErrorResponse(c, err)
		return
	}
	c.JSON(http.StatusOK, liveVo)
}

func (m memoryController) Days(c *gin.Context) {
	days, version, err := m.service.Days()
	if err != nil {
		util.ErrorResponse(c, err)
		return
	}
	etag := strconv.FormatInt(version, 10)
	c.Header("ETag", etag)
	c.Header("Cache-Control", "public, max-age=300, must-revalidate")
	if match := c.GetHeader("If-None-Match"); match == etag {
		c.AbortWithStatus(http.StatusNotModified)
		return
	}
	c.JSON(http.StatusOK, days)
}

func (m memoryController) Day(c *gin.Context) {
	day := c.Param("day")
	data, version, err := m.service.Day(day)
	if err != nil {
		util.ErrorResponse(c, err)
		return
	}
	etag := strconv.FormatInt(version, 10)
	c.Header("ETag", etag)
	if c.Query("v") != "" {
		c.Header("Cache-Control", "public, max-age=31536000, must-revalidate")
	} else {
		c.Header("Cache-Control", "public, max-age=300, must-revalidate")
	}
	if match := c.GetHeader("If-None-Match"); match == etag {
		c.AbortWithStatus(http.StatusNotModified)
		return
	}
	c.JSON(http.StatusOK, data)
}
