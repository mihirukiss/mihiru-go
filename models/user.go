package models

type UserBaseFields struct {
	Name      string   `bson:"name" json:"name"`
	LoginName string   `bson:"loginName" json:"loginName"`
	Roles     []string `bson:"roles" json:"roles"`
}

type UserSecurityField struct {
	Password string `bson:"password" json:"password"`
}

type User struct {
	UserBaseFields    `bson:",inline"`
	UserSecurityField `bson:",inline"`
}

type UserWithObjectId struct {
	ObjectIdFields `bson:",inline"`
	User           `bson:",inline"`
}
