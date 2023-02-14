package utils

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"math/rand"
)

func GetRandom(seed ...interface{}) int64 {
	str := primitive.NewObjectID().Hex()
	for _, v := range seed {
		str += AsString(v)
	}
	return rand.New(rand.NewSource(GetHashCode(str))).Int63()
}

func GetRandomN(n int64, seed ...interface{}) int64 {
	str := primitive.NewObjectID().Hex()
	for _, v := range seed {
		str += AsString(v)
	}
	return rand.New(rand.NewSource(GetHashCode(str))).Int63n(n)
}
