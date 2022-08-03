package mongo

import "go.mongodb.org/mongo-driver/bson/primitive"

type UID string

func NewUID() string {
	objid := primitive.NewObjectID()
	str := objid.String()
	return str
}
