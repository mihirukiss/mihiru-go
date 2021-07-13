package models

type DynamicUserProfile struct {
	Uid   int64  `bson:"uid" json:"uid"`
	Uname string `bson:"uname" json:"uname"`
	Face  string `bson:"face" json:"face"`
	Utype string `bson:"utype" json:"utype"`
}

type DynamicPicture struct {
	ImgSrc    string `bson:"img_src" json:"img_src"`
	ImgWidth  int    `bson:"img_width" json:"img_width"`
	ImgHeight int    `bson:"img_height" json:"img_height"`
	ImgSize   int    `bson:"img_size" json:"img_size"`
}

type DynamicCtrl struct {
	Data     string `bson:"data" json:"data"`
	Length   int    `bson:"length" json:"length"`
	Location int    `bson:"location" json:"location"`
	Type     int16  `bson:"type" json:"type"`
}

type Dynamic struct {
	Type         int16                 `bson:"type" json:"type"`
	DynamicId    int64                 `bson:"dynamic_id" json:"dynamic_id"`
	Rid          int64                 `bson:"rid" json:"rid"`
	Bvid         string                `bson:"bvid" json:"bvid"`
	Timestamp    int64                 `bson:"timestamp" json:"timestamp"`
	UserProfile  *DynamicUserProfile   `bson:"user_profile" json:"user_profile"`
	UserProfiles *[]DynamicUserProfile `bson:"user_profiles" json:"user_profiles"`
	Title        string                `bson:"title" json:"title"`
	Content      string                `bson:"content" json:"content"` //content, summary, description, dynamic
	Desc         string                `bson:"desc" json:"desc"`
	Pic          string                `bson:"pic" json:"pic"` //pic, cover
	AreaV2Name   string                `bson:"area_v2_name" json:"area_v2_name"`
	ImageUrls    []string              `bson:"image_urls" json:"image_urls"`
	Ctrl         *[]DynamicCtrl        `bson:"ctrl" json:"ctrl"`
	Pictures     *[]DynamicPicture     `bson:"pictures" json:"pictures"`
	Origin       *Dynamic              `bson:"origin" json:"origin"`
}

type DynamicWithLastModified struct {
	Dynamic            `bson:",inline"`
	LastModifiedFields `bson:",inline"`
}

type DynamicWithObjectId struct {
	ObjectIdFields          `bson:",inline"`
	DynamicWithLastModified `bson:",inline"`
}
