package services

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"mihiru-go/database"
	"mihiru-go/models"
	"mihiru-go/util"
	"mihiru-go/vo"
	"net/http"
	"time"
)

type MemoryService interface {
	AddDynamic(dynamic *models.Dynamic) (*models.DynamicWithObjectId, error)
	UpdateDynamic(id primitive.ObjectID, dynamic *models.Dynamic) (*models.DynamicWithObjectId, error)
	AddLive(live *models.Live) (*models.LiveWithObjectId, error)
	UpdateLive(id primitive.ObjectID, live *models.Live) (*models.LiveWithObjectId, error)
	Days() ([]*models.DayCount, int64, error)
	Day(day string) ([]interface{}, int64, error)
}

type memoryService struct {
	dynamicDatabase database.DynamicDatabase
	liveDatabase    database.LiveDatabase
}

var dayCountsCache []*models.DayCount
var dayCountsVersionCache int64
var dayCacheMap = make(map[string][]interface{})
var dayVersionCacheMap = make(map[string]int64)

func NewMemoryService(dynamicDatabase database.DynamicDatabase, liveDatabase database.LiveDatabase) MemoryService {
	return memoryService{dynamicDatabase, liveDatabase}
}

func (m memoryService) AddDynamic(dynamic *models.Dynamic) (*models.DynamicWithObjectId, error) {
	dynamicWithObjectId := new(models.DynamicWithObjectId)
	dynamicWithObjectId.Dynamic = *dynamic
	dynamicWithObjectId.LastModified = time.Now().UnixNano() / 1e6
	err := m.dynamicDatabase.InsertDynamic(dynamicWithObjectId)
	if err != nil {
		util.LogError(err)
		return nil, vo.NewErrorWithHttpStatus("添加数据失败, 请稍后重试", http.StatusInternalServerError)
	}
	cleanCache(dynamic.Timestamp)
	return dynamicWithObjectId, nil
}

func (m memoryService) UpdateDynamic(id primitive.ObjectID, dynamic *models.Dynamic) (*models.DynamicWithObjectId, error) {
	dynamicWithObjectId := new(models.DynamicWithObjectId)
	dynamicWithObjectId.ID = id
	dynamicWithObjectId.Dynamic = *dynamic
	dynamicWithObjectId.LastModified = time.Now().UnixNano() / 1e6
	err := m.dynamicDatabase.UpdateDynamic(dynamicWithObjectId)
	if err != nil {
		util.LogError(err)
		return nil, vo.NewErrorWithHttpStatus("更新数据失败, 请稍后重试", http.StatusInternalServerError)
	}
	cleanCache(dynamic.Timestamp)
	return dynamicWithObjectId, nil
}

func (m memoryService) AddLive(live *models.Live) (*models.LiveWithObjectId, error) {
	liveWithObjectId := new(models.LiveWithObjectId)
	liveWithObjectId.Live = *live
	liveWithObjectId.LastModified = time.Now().UnixNano() / 1e6
	err := m.liveDatabase.InsertLive(liveWithObjectId)
	if err != nil {
		util.LogError(err)
		return nil, vo.NewErrorWithHttpStatus("添加数据失败, 请稍后重试", http.StatusInternalServerError)
	}
	cleanCache(live.Timestamp)
	return liveWithObjectId, nil
}

func (m memoryService) UpdateLive(id primitive.ObjectID, live *models.Live) (*models.LiveWithObjectId, error) {
	liveWithObjectId := new(models.LiveWithObjectId)
	liveWithObjectId.ID = id
	liveWithObjectId.Live = *live
	liveWithObjectId.LastModified = time.Now().UnixNano() / 1e6
	err := m.liveDatabase.UpdateLive(liveWithObjectId)
	if err != nil {
		util.LogError(err)
		return nil, vo.NewErrorWithHttpStatus("更新数据失败, 请稍后重试", http.StatusInternalServerError)
	}
	cleanCache(live.Timestamp)
	return liveWithObjectId, nil
}

func (m memoryService) Days() ([]*models.DayCount, int64, error) {
	if dayCountsCache == nil {
		dynamicCount, err := m.dynamicDatabase.CountDynamicByDay()
		if err != nil {
			util.LogError(err)
			return nil, 0, vo.NewErrorWithHttpStatus("统计动态数据失败, 请稍后重试", http.StatusInternalServerError)
		}
		liveCount, err := m.liveDatabase.CountLiveByDay()
		if err != nil {
			util.LogError(err)
			return nil, 0, vo.NewErrorWithHttpStatus("统计直播数据失败, 请稍后重试", http.StatusInternalServerError)
		}
		var data []*models.DayCount
		dynamicIndex := 0
		liveIndex := 0
		maxVersion := int64(0)
		dynamicLen := len(dynamicCount)
		liveLen := len(liveCount)
		for dynamicIndex < dynamicLen || liveIndex < liveLen {
			if dynamicIndex >= dynamicLen {
				data = append(data, liveCount[liveIndex])
				if liveCount[liveIndex].Version > maxVersion {
					maxVersion = liveCount[liveIndex].Version
				}
				liveIndex++
			} else if liveIndex >= liveLen {
				data = append(data, dynamicCount[dynamicIndex])
				if dynamicCount[dynamicIndex].Version > maxVersion {
					maxVersion = dynamicCount[dynamicIndex].Version
				}
				dynamicIndex++
			} else if liveCount[liveIndex].Day == dynamicCount[dynamicIndex].Day {
				mergeDayCount := new(models.DayCount)
				mergeDayCount.Day = liveCount[liveIndex].Day
				mergeDayCount.Count = liveCount[liveIndex].Count + dynamicCount[dynamicIndex].Count
				if liveCount[liveIndex].Version > dynamicCount[dynamicIndex].Version {
					mergeDayCount.Version = liveCount[liveIndex].Version
				} else {
					mergeDayCount.Version = dynamicCount[dynamicIndex].Version
				}
				data = append(data, mergeDayCount)
				if mergeDayCount.Version > maxVersion {
					maxVersion = mergeDayCount.Version
				}
				liveIndex++
				dynamicIndex++
			} else if liveCount[liveIndex].Day > dynamicCount[dynamicIndex].Day {
				data = append(data, dynamicCount[dynamicIndex])
				if dynamicCount[dynamicIndex].Version > maxVersion {
					maxVersion = dynamicCount[dynamicIndex].Version
				}
				dynamicIndex++
			} else {
				data = append(data, liveCount[liveIndex])
				if liveCount[liveIndex].Version > maxVersion {
					maxVersion = liveCount[liveIndex].Version
				}
				liveIndex++
			}
		}
		dayCountsCache = data
		dayCountsVersionCache = maxVersion
	}
	return dayCountsCache, dayCountsVersionCache, nil
}

func (m memoryService) Day(day string) ([]interface{}, int64, error) {
	result := dayCacheMap[day]
	if result != nil {
		return result, dayVersionCacheMap[day], nil
	}
	result = []interface{}{}
	loc, _ := time.LoadLocation("Asia/Shanghai")
	date, err := time.ParseInLocation("2006.01.02", day, loc)
	if err != nil {
		return nil, 0, vo.NewErrorWithHttpStatus("无效的日期参数", http.StatusBadRequest)
	}
	startTimestamp := date.Unix()
	date = date.AddDate(0, 0, 1)
	endTimestamp := date.Unix()
	dynamics, err := m.dynamicDatabase.QueryDynamicByTimestamp(startTimestamp, endTimestamp)
	if err != nil {
		util.LogError(err)
		return nil, 0, vo.NewErrorWithHttpStatus("查询动态数据失败, 请稍后重试", http.StatusInternalServerError)
	}
	lives, err := m.liveDatabase.QueryLiveByTimestamp(startTimestamp, endTimestamp)
	if err != nil {
		util.LogError(err)
		return nil, 0, vo.NewErrorWithHttpStatus("查询直播数据失败, 请稍后重试", http.StatusInternalServerError)
	}
	var maxVersion int64
	dynamicIndex := 0
	liveIndex := 0
	dynamicLen := len(dynamics)
	liveLen := len(lives)
	for dynamicIndex < dynamicLen || liveIndex < liveLen {
		if dynamicIndex >= dynamicLen {
			result = append(result, lives[liveIndex])
			if lives[liveIndex].LastModified > maxVersion {
				maxVersion = lives[liveIndex].LastModified
			}
			liveIndex++
		} else if liveIndex >= liveLen {
			result = append(result, dynamics[dynamicIndex])
			if dynamics[dynamicIndex].LastModified > maxVersion {
				maxVersion = dynamics[dynamicIndex].LastModified
			}
			dynamicIndex++
		} else if dynamics[dynamicIndex].Timestamp > lives[liveIndex].Timestamp {
			result = append(result, lives[liveIndex])
			if lives[liveIndex].LastModified > maxVersion {
				maxVersion = lives[liveIndex].LastModified
			}
			liveIndex++
		} else {
			result = append(result, dynamics[dynamicIndex])
			if dynamics[dynamicIndex].LastModified > maxVersion {
				maxVersion = dynamics[dynamicIndex].LastModified
			}
			dynamicIndex++
		}
	}
	dayCacheMap[day] = result
	dayVersionCacheMap[day] = maxVersion
	return result, maxVersion, nil
}

func cleanCache(timestamp int64) {
	dayCountsCache = nil
	loc, _ := time.LoadLocation("Asia/Shanghai")
	delete(dayCacheMap, time.Unix(timestamp, 0).In(loc).Format("2006.01.02"))
}
