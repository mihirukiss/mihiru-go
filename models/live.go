package models

type LiveCut struct {
	Type      int8               `bson:"type" json:"type"`
	File      string             `bson:"file" json:"file"`
	Title     string             `bson:"title" json:"title"`
	Dialogues *[]LiveCutDialogue `bson:"dialogues" json:"dialogues"`
}

type LiveCutDialogue struct {
	Speaker string `bson:"speaker" json:"speaker"`
	Text    string `bson:"text" json:"text"`
}

type Live struct {
	Timestamp        int64                 `bson:"timestamp" json:"timestamp"`
	Title            string                `bson:"title" json:"title"`
	Content          string                `bson:"content" json:"content"`
	UserProfile      *DynamicUserProfile   `bson:"user_profile" json:"user_profile"`
	JoinUserProfiles *[]DynamicUserProfile `bson:"join_user_profiles" json:"join_user_profiles"`
	FullRecord       string                `bson:"full_record" json:"full_record"`
	Cuts             *[]LiveCut            `bson:"cuts" json:"cuts"`
}

type LiveWithLastModified struct {
	Live               `bson:",inline"`
	LastModifiedFields `bson:",inline"`
}

type LiveWithObjectId struct {
	ObjectIdFields       `bson:",inline"`
	LiveWithLastModified `bson:",inline"`
}
