package vo

import "mihiru-go/models"

type LiveVoicesCategoryVo struct {
	Category string    `json:"category"`
	Voices   []VoiceVo `json:"voices"`
}

type VoiceVo struct {
	models.ObjectIdFields
	models.VoiceBaseFields
	models.VoiceAutoGenFields
}
