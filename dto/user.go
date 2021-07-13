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

type ChangePasswordDto struct {
	OldPassword string `json:"old_password"`
	NewPassword string `json:"new_password"`
}
