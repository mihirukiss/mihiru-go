package services

import (
	"crypto/hmac"
	"crypto/sha512"
	"encoding/hex"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"log"
	"math/rand"
	"mihiru-go/config"
	"mihiru-go/database"
	"mihiru-go/dto"
	"mihiru-go/models"
	"mihiru-go/util"
	"mihiru-go/vo"
	"net/http"
	"time"
	"unsafe"
)

type UserService interface {
	Add(userDto *dto.UserDto) (*vo.UserVo, error)
	Login(loginDto *dto.LoginDto) (string, string, error)
	ChangePassword(token string, changePasswordDto dto.ChangePasswordDto) error
	CheckToken(token string) *vo.UserVo
	InitUser()
}

type userService struct {
	passwordEncoderKey []byte
	db                 database.UserDatabase
}

var tokenMap = make(map[primitive.ObjectID]string)
var loginInfoMap = make(map[string]*vo.UserVo)
var randSource = rand.NewSource(time.Now().UnixNano())

const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
const (
	letterIdxBits = 6                    // 6 bits to represent a letter index
	letterIdxMask = 1<<letterIdxBits - 1 // All 1-bits, as many as letterIdxBits
	letterIdxMax  = 63 / letterIdxBits   // # of letter indices fitting in 63 bits
)

func NewUserService(db database.UserDatabase) UserService {
	return userService{
		passwordEncoderKey: []byte(config.GetConfigs().GetString("security.password-key")),
		db:                 db,
	}
}

func (u userService) Add(userDto *dto.UserDto) (*vo.UserVo, error) {
	user := new(models.UserWithObjectId)
	user.LoginName = userDto.LoginName
	user.Password = encodePassword(userDto.Password, u.passwordEncoderKey)
	user.Name = userDto.Name
	user.Roles = userDto.Roles
	err := u.db.InsertUser(user)
	if err != nil {
		util.LogError(err)
		return nil, vo.NewErrorWithHttpStatus("添加数据失败, 请稍后重试", http.StatusInternalServerError)
	}
	return convertToUserVo(user), nil
}

func (u userService) ChangePassword(token string, changePasswordDto dto.ChangePasswordDto) error {
	userVo := loginInfoMap[token]
	if userVo == nil {
		return vo.NewErrorWithHttpStatus("用户未登录或登录已失效", http.StatusForbidden)
	}
	userWithObjectId, err := u.db.GetUserById(userVo.ID)
	if err != nil {
		util.LogError(err)
		return vo.NewErrorWithHttpStatus("获取登录用户信息失败, 请稍后重试", http.StatusInternalServerError)
	}
	if userWithObjectId == nil {
		return vo.NewErrorWithHttpStatus("无法获取当前用户信息, 请稍后重试", http.StatusInternalServerError)
	}
	if encodePassword(changePasswordDto.OldPassword, u.passwordEncoderKey) != userWithObjectId.Password {
		return vo.NewErrorWithHttpStatus("原密码校验失败", http.StatusBadRequest)
	}
	userWithObjectId.Password = encodePassword(changePasswordDto.NewPassword, u.passwordEncoderKey)
	err = u.db.UpdateUser(userWithObjectId)
	if err != nil {
		return vo.NewErrorWithHttpStatus("更新用户信息失败, 请稍后重试", http.StatusInternalServerError)
	}
	delete(loginInfoMap, token)
	return nil
}

func (u userService) Login(loginDto *dto.LoginDto) (string, string, error) {
	user, err := u.db.GetUserByLoginName(loginDto.LoginName)
	if err != nil {
		util.LogError(err)
		return "", "", vo.NewErrorWithHttpStatus("查询用户信息失败, 请稍后重试", http.StatusInternalServerError)
	}
	if user == nil || encodePassword(loginDto.Password, u.passwordEncoderKey) != user.Password {
		return "", "", vo.NewErrorWithHttpStatus("账号或密码错误", http.StatusBadRequest)
	}
	token := randString(64, randSource)
	loginInfoMap[token] = convertToUserVo(user)
	delete(loginInfoMap, tokenMap[user.ID])
	tokenMap[user.ID] = token
	return token, user.Name, nil
}

func (u userService) CheckToken(token string) *vo.UserVo {
	return loginInfoMap[token]
}

func (u userService) InitUser() {
	initLoginName := config.GetConfigs().GetString("security.init-login-name")
	user, err := u.db.GetUserByLoginName(initLoginName)
	if err != nil {
		log.Fatal(err.Error())
	} else if user == nil {
		user = new(models.UserWithObjectId)
		user.Name = initLoginName
		user.LoginName = initLoginName
		user.Password = encodePassword(config.GetConfigs().GetString("security.init-password"), u.passwordEncoderKey)
		user.Roles = []string{"admin"}
		err = u.db.InsertUser(user)
		if err != nil {
			util.LogError(err)
		}
	}
}

func encodePassword(password string, passwordEncoderKey []byte) string {
	hash := hmac.New(sha512.New, passwordEncoderKey)
	hash.Write([]byte(password))
	return hex.EncodeToString(hash.Sum(nil))
}

func convertToUserVo(user *models.UserWithObjectId) *vo.UserVo {
	userVo := new(vo.UserVo)
	userVo.ID = user.ID
	userVo.UserBaseFields = user.UserBaseFields
	return userVo
}

func randString(n int, src rand.Source) string {
	b := make([]byte, n)
	// A src.Int63() generates 63 random bits, enough for letterIdxMax characters!
	for i, cache, remain := n-1, src.Int63(), letterIdxMax; i >= 0; {
		if remain == 0 {
			cache, remain = src.Int63(), letterIdxMax
		}
		if idx := int(cache & letterIdxMask); idx < len(letterBytes) {
			b[i] = letterBytes[idx]
			i--
		}
		cache >>= letterIdxBits
		remain--
	}

	return *(*string)(unsafe.Pointer(&b))
}
