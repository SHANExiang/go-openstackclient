package utils

import (
	"strings"
	"unicode"
)

// Pluralize 函数接受一个字符串作为参数，返回该字符串的复数形式。
func Pluralize(word string) string {
	// 如果输入字符串为空，返回默认值或给出错误提示。
	if word == "" {
		return "" // 或者返回默认值或错误提示语句
	}

	// 将字符串转为小写处理。
	word = strings.ToLower(word)

	// 根据不同规则进行复数形式的转换。
	switch {
	case strings.HasSuffix(word, "s") || strings.HasSuffix(word, "x") || strings.HasSuffix(word, "z"):
		return word + "es"
	case strings.HasSuffix(word, "y") && len(word) > 1 && !isVowel(word[len(word)-2]):
		return word[:len(word)-1] + "ies"
	default:
		return word + "s"
	}
}

// 判断字符 c 是否为元音字母
func isVowel(c byte) bool {
	s := unicode.ToLower(rune(c))
	return s == 'a' || c == 'e' || c == 'i' || c == 'o' || c == 'u'
}
