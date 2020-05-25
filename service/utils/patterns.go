package utils

import (
	"regexp"
	"strings"
)

// 进行路径匹配
// 路径分隔符为 '/'，可在 pattern 结尾加上 '**' 进行匹配
//
func URLPathMatch(pattern, path string) bool {

	patterns := strings.Split(pattern, "/")
	paths := strings.Split(path, "/")

	patternIndexStart := 0
	patternIndexEnd := len(patterns) - 1
	pathIndexStart := 0
	pathIndexEnd := len(paths) - 1

	for patternIndexEnd >= patternIndexStart {
		patt := patterns[patternIndexStart]
		if patt == "**" {
			break
		}
		if pathIndexEnd < pathIndexStart || patt != paths[pathIndexStart] {
			return false
		}
		patternIndexStart++
		pathIndexStart++
	}

	if patternIndexEnd < patternIndexStart && pathIndexStart <= pathIndexEnd {
		return false
	}

	if patternIndexStart != patternIndexEnd {
		return false
	}

	return true
}


func TrimHtml(src string) string {
	//将HTML标签全转换成小写
	re, _ := regexp.Compile("\\<[\\S\\s]+?\\>")
	src = re.ReplaceAllStringFunc(src, strings.ToLower)
	//去除STYLE
	re, _ = regexp.Compile("\\<style[\\S\\s]+?\\</style\\>")
	src = re.ReplaceAllString(src, "")
	//去除SCRIPT
	re, _ = regexp.Compile("\\<script[\\S\\s]+?\\</script\\>")
	src = re.ReplaceAllString(src, "")
	//去除所有尖括号内的HTML代码，并换成换行符
	re, _ = regexp.Compile("\\<[\\S\\s]+?\\>")
	src = re.ReplaceAllString(src, "\n")
	//去除连续的换行符
	re, _ = regexp.Compile("\\s{2,}")
	src = re.ReplaceAllString(src, "\n")
	return strings.TrimSpace(src)
}