package utils

import "strings"

// UnderLineToCamel 将下划线风格的单词变为驼峰命名的单词
func UnderLineToCamel(line string) string {
	words := strings.Split(line, "_")
	n := ""
	for _, w := range words {
		n += strings.ToUpper(w[0:1]) + w[1:]
	}
	return n
}

// SnakeString 驼峰转蛇形
func SnakeString(s string) string {
	data := make([]byte, 0, len(s)*2)
	j := false
	num := len(s)
	for i := 0; i < num; i++ {
		d := s[i]
		if i > 0 && d >= 'A' && d <= 'Z' && j {
			data = append(data, '_')
		}
		if d != '_' {
			j = true
		}
		data = append(data, d)
	}
	return strings.ToLower(string(data[:]))
}
