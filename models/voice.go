package models

type VoiceBaseFields struct {
	Liver    string `bson:"liver" json:"liver"`
	Category string `bson:"category" json:"category"`
	Title    string `bson:"title" json:"title"`
	SortNo   int64  `bson:"sort_no" json:"sort_no"`
	Remark   string `bson:"remark" json:"remark"`
}

type VoiceAutoGenFields struct {
	FilePath string `bson:"file_path" json:"file_path"`
	AddTime  int64  `bson:"add_time" json:"add_time"`
}

type Voice struct {
	VoiceBaseFields    `bson:",inline"`
	VoiceAutoGenFields `bson:",inline"`
	Deleted            bool `bson:"deleted" json:"deleted"`
}

type VoiceWithObjectId struct {
	ObjectIdFields `bson:",inline"`
	Voice          `bson:",inline"`
}
