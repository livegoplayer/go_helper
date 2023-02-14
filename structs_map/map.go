package structs_map

import (
	"github.com/fatih/structs"
)

func NewStructMap(st interface{}, tagName string) *structs.Struct {
	stt := structs.New(st)
	if tagName != "" {
		stt.TagName = tagName
	}
	return stt
}
