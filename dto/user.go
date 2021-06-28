package dto

import "mihiru-go/models"

type UserDto struct {
	models.UserBaseFields
	models.UserSecurityField
}

type LoginDto struct {
	LoginName string `json:"loginName"`
	Password  string `json:"password"`
}
