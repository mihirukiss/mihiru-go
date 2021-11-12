package services

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"mihiru-go/config"
	"mihiru-go/database"
	"mihiru-go/models"
	"mihiru-go/util"
	"mihiru-go/vo"
	"mime/multipart"
	"net/http"
	"os"
	"time"
)

type VoiceService interface {
	AddVoice(voiceDto *models.VoiceBaseFields, file *multipart.FileHeader, c *gin.Context) error
	UpdateVoice(id primitive.ObjectID, voiceDto *models.VoiceBaseFields, file *multipart.FileHeader, c *gin.Context) error
	DeleteVoice(id primitive.ObjectID) error
	LiverVoices(liver string) ([]*vo.LiveVoicesCategoryVo, int64, error)
}

type voiceService struct {
	voiceDatabase database.VoiceDatabase
}

var liveVoiceCacheMap = make(map[string][]*vo.LiveVoicesCategoryVo)
var liveVoiceVersionCacheMap = make(map[string]int64)

func NewVoiceService(voiceDatabase database.VoiceDatabase) VoiceService {
	return voiceService{voiceDatabase}
}

func (v voiceService) AddVoice(voiceDto *models.VoiceBaseFields, file *multipart.FileHeader, c *gin.Context) error {
	voice := new(models.Voice)
	voice.VoiceBaseFields = *voiceDto
	voice.AddTime = time.Now().UnixNano() / 1e6
	if voice.SortNo <= 0 {
		voice.SortNo = voice.AddTime
	}
	fileName := uuid.New().String() + ".mp3"
	configs := config.GetConfigs()
	folderPath := configs.GetString("voice.base-folder") + voice.Liver
	filePath := folderPath + "/" + fileName
	urlPath := configs.GetString("voice.base-path") + voice.Liver + "/" + fileName
	voice.FilePath = urlPath
	voice.Deleted = false
	err := createFolderIfNotExists(folderPath)
	if err != nil {
		util.LogError(err)
		return vo.NewErrorWithHttpStatus("建立保存文件夹失败, 请稍后重试", http.StatusInternalServerError)
	}
	err = c.SaveUploadedFile(file, filePath)
	if err != nil {
		util.LogError(err)
		return vo.NewErrorWithHttpStatus("保存文件失败, 请稍后重试", http.StatusInternalServerError)
	}
	err = v.voiceDatabase.InsertVoice(voice)
	if err != nil {
		util.LogError(err)
		return vo.NewErrorWithHttpStatus("添加数据失败, 请稍后重试", http.StatusInternalServerError)
	}
	_, _, err = v.refreshCache(voice.Liver)
	if err != nil {
		util.LogError(err)
		return vo.NewErrorWithHttpStatus("刷新失败, 请稍后重试", http.StatusInternalServerError)
	}
	return nil
}

func (v voiceService) UpdateVoice(id primitive.ObjectID, voiceDto *models.VoiceBaseFields, file *multipart.FileHeader, c *gin.Context) error {
	voice, err := v.voiceDatabase.GetVoiceById(id)
	if err != nil {
		util.LogError(err)
		return vo.NewErrorWithHttpStatus("查找数据失败, 请稍后重试", http.StatusInternalServerError)
	}
	if voice == nil {
		return vo.NewErrorWithHttpStatus("数据不存在", http.StatusNotFound)
	}
	voice.VoiceBaseFields = *voiceDto
	if voice.SortNo <= 0 {
		voice.SortNo = voice.AddTime
	}
	if file != nil {
		fileName := uuid.New().String() + ".mp3"
		configs := config.GetConfigs()
		folderPath := configs.GetString("voice.base-folder") + voice.Liver
		filePath := folderPath + "/" + fileName
		urlPath := configs.GetString("voice.base-path") + voice.Liver + "/" + fileName
		voice.FilePath = urlPath
		err = createFolderIfNotExists(folderPath)
		if err != nil {
			util.LogError(err)
			return vo.NewErrorWithHttpStatus("建立保存文件夹失败, 请稍后重试", http.StatusInternalServerError)
		}
		err = c.SaveUploadedFile(file, filePath)
		if err != nil {
			util.LogError(err)
			return vo.NewErrorWithHttpStatus("保存文件失败, 请稍后重试", http.StatusInternalServerError)
		}
	}
	err = v.voiceDatabase.UpdateVoice(voice)
	if err != nil {
		util.LogError(err)
		return vo.NewErrorWithHttpStatus("添加数据失败, 请稍后重试", http.StatusInternalServerError)
	}
	_, _, err = v.refreshCache(voice.Liver)
	if err != nil {
		util.LogError(err)
		return vo.NewErrorWithHttpStatus("刷新失败, 请稍后重试", http.StatusInternalServerError)
	}
	return nil
}

func (v voiceService) DeleteVoice(id primitive.ObjectID) error {
	voice, err := v.voiceDatabase.GetVoiceById(id)
	if err != nil {
		util.LogError(err)
		return vo.NewErrorWithHttpStatus("查找数据失败, 请稍后重试", http.StatusInternalServerError)
	}
	if voice == nil {
		return vo.NewErrorWithHttpStatus("数据不存在", http.StatusNotFound)
	}
	err = v.voiceDatabase.DeleteVoice(id)
	if err != nil {
		util.LogError(err)
		return vo.NewErrorWithHttpStatus("删除数据失败, 请稍后重试", http.StatusInternalServerError)
	}
	_, _, err = v.refreshCache(voice.Liver)
	if err != nil {
		util.LogError(err)
		return vo.NewErrorWithHttpStatus("刷新失败, 请稍后重试", http.StatusInternalServerError)
	}
	return nil
}

func (v voiceService) LiverVoices(liver string) ([]*vo.LiveVoicesCategoryVo, int64, error) {
	result := liveVoiceCacheMap[liver]
	if result != nil {
		return result, liveVoiceVersionCacheMap[liver], nil
	}
	return v.refreshCache(liver)
}

func (v voiceService) refreshCache(liver string) ([]*vo.LiveVoicesCategoryVo, int64, error) {
	categories := make(map[string]*vo.LiveVoicesCategoryVo)
	var result []*vo.LiveVoicesCategoryVo
	voices, err := v.voiceDatabase.ListVoiceByLiver(liver)
	if err != nil {
		util.LogError(err)
		return nil, 0, vo.NewErrorWithHttpStatus("查询数据失败, 请稍后重试", http.StatusInternalServerError)
	}
	if voices != nil && len(voices) > 0 {
		for _, voice := range voices {
			if _, existed := categories[voice.Category]; !existed {
				category := new(vo.LiveVoicesCategoryVo)
				category.Category = voice.Category
				categories[voice.Category] = category
				category.Voices = []vo.VoiceVo{}
				result = append(result, category)
			}
			voiceVo := new(vo.VoiceVo)
			voiceVo.ObjectIdFields = voice.ObjectIdFields
			voiceVo.VoiceBaseFields = voice.VoiceBaseFields
			voiceVo.VoiceAutoGenFields = voice.VoiceAutoGenFields
			categories[voice.Category].Voices = append(categories[voice.Category].Voices, *voiceVo)
		}
	}
	liveVoiceCacheMap[liver] = result
	liveVoiceVersionCacheMap[liver] = time.Now().UnixNano() / 1e6
	return result, liveVoiceVersionCacheMap[liver], nil
}

func createFolderIfNotExists(folder string) error {
	_, err := os.Stat(folder)
	if err != nil && os.IsNotExist(err) {
		return os.MkdirAll(folder, os.ModePerm)
	}
	return nil
}
