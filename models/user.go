package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type UserBaseFields struct {
	Name      string   `bson:"name" json:"name"`
	LoginName string   `bson:"loginName" json:"loginName"`
	Roles     []string `bson:"roles" json:"roles"`
}

type UserSecurityField struct {
	Password string `bson:"password" json:"password"`
}

type UserObjectIdFields struct {
	ID primitive.ObjectID `bson:"_id" json:"id"`
}

type User struct {
	UserBaseFields    `bson:",inline"`
	UserSecurityField `bson:",inline"`
}

type UserWithObjectId struct {
	UserObjectIdFields `bson:",inline"`
	User               `bson:",inline"`
}
