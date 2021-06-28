package util

import (
	"github.com/gin-gonic/gin"
	"log"
	"mihiru-go/vo"
	"net/http"
	"runtime"
	"sort"
)

func TextInArray(target string, sortedArray []string) bool {
	index := sort.SearchStrings(sortedArray, target)
	if index < len(sortedArray) && sortedArray[index] == target {
		return true
	}
	return false
}

func LogError(err error) {
	if err != nil {
		pc, fn, line, _ := runtime.Caller(1)
		log.Printf("[error] in %s[%s:%d] %v", runtime.FuncForPC(pc).Name(), fn, line, err)
	}
}

func ErrorResponse(c *gin.Context, err error) {
	if e, ok := err.(vo.ErrorWithHttpStatus); ok {
		c.AbortWithStatusJSON(e.HttpStatus(), gin.H{"message": e.Error()})
		return
	}
	LogError(err)
	c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"message": "出现未知错误, 请稍后重试"})
}
