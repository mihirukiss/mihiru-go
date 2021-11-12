package controllers

import (
	"github.com/gin-gonic/gin"
	"log"
	"mihiru-go/dto"
	"mihiru-go/services"
	"mihiru-go/util"
	"net/http"
)

type UserController interface {
	Add(c *gin.Context)
	Login(c *gin.Context)
	ChangePassword(c *gin.Context)
}

type userController struct {
	service services.UserService
}

func NewUserController(service services.UserService) UserController {
	return userController{service: service}
}

func (u userController) Add(c *gin.Context) {
	var userDto dto.UserDto
	if err := c.BindJSON(&userDto); err != nil {
		log.Println(err.Error())
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": "错误的参数格式"})
		return
	}
	userVo, err := u.service.Add(&userDto)
	if err != nil {
		util.ErrorResponse(c, err)
		return
	}
	c.JSON(http.StatusOK, userVo)
}

func (u userController) ChangePassword(c *gin.Context) {
	var changePasswordDto dto.ChangePasswordDto
	if err := c.BindJSON(&changePasswordDto); err != nil {
		log.Println(err.Error())
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": "错误的参数格式"})
		return
	}
	err := u.service.ChangePassword(c.GetHeader("authorization"), changePasswordDto)
	if err != nil {
		util.ErrorResponse(c, err)
		return
	}
	c.Status(http.StatusOK)
}

func (u userController) Login(c *gin.Context) {
	var loginDto dto.LoginDto
	if err := c.BindJSON(&loginDto); err != nil {
		log.Println(err.Error())
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": "错误的参数格式"})
		return
	}
	token, name, err := u.service.Login(&loginDto)
	if err != nil {
		util.ErrorResponse(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"token": token, "name": name})
}
