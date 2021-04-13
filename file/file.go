package file

import (
	"regexp"
	"strings"
)

const (
	FOLDER = 1
	IMG    = 2
	VIDEO  = 3
	OTHER  = 4
	PDF    = 5
	OFFICE = 6
	TXT    = 7
)

func GetFileTypeByName(filename string) int {
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
