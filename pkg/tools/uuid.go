package tools

import "go.mongodb.org/mongo-driver/bson/primitive"

type UUID string

func NewUUID() UUID {
	objid := primitive.NewObjectID()
	str := objid.String()
	return UUID(str)
}

func (id UUID) String() string {
	return string(id)
}
