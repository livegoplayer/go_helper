package utils

import (
	"strconv"
	"strings"
)

func VersionThan(va string, vb string) int {
	vas := strings.Split(va, ".")
	vbs := strings.Split(vb, ".")

	count := len(vas)
	if len(vas) < len(vbs) {
		count = len(vbs)
	}
	for i := 0; i < count; i++ {
		a := 0
		if len(vas) < i+1 {
			a = 0
		} else {
			a, _ = strconv.Atoi(vas[i])
		}

		b := 0
		if len(vbs) < i+1 {
			b = 0
		} else {
			b, _ = strconv.Atoi(vbs[i])
		}
		if a > b {
			return 1
		}
		if a < b {
			return -1
		}
	}
	return 0
}
