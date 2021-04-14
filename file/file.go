package file

import (
	"regexp"
	"strings"
)

type FileType int64

const (
	FOLDER FileType = 1
	IMG    FileType = 2
	VIDEO  FileType = 3
	OTHER  FileType = 4
	PDF    FileType = 5
	OFFICE FileType = 6
	TXT    FileType = 7
)

func GetFileTypeByName(filename string) FileType {
	if match, _ := regexp.MatchString("(.*)\\.(jpg|bmp|gif|ico|pcx|jpeg|tif|png|raw|tga)", strings.ToLower(filename)); match {
		return IMG
	}

	if match, _ := regexp.MatchString("(.*)\\.(swf|flv|mp4|rmvb|avi|mpeg|ra|ram|mov|wmv)", strings.ToLower(filename)); match {
		return VIDEO
	}

	if match, _ := regexp.MatchString("(.*)\\.(pdf)", strings.ToLower(filename)); match {
		return PDF
	}

	if match, _ := regexp.MatchString("(.*)\\.(docx|dotx|xlsx|xltx|pptx|potx|ppsx)", strings.ToLower(filename)); match {
		return OFFICE
	}

	if match, _ := regexp.MatchString("(.*)\\.(txt|me|sql|json|log)", strings.ToLower(filename)); match {
		return TXT
	}

	return OTHER
}
