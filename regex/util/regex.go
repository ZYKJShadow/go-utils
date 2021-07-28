package util

import (
	"errors"
	"regexp"
	"strconv"
)

func MustFindAllString(s string, reg string) (res []string) {
	compile := regexp.MustCompile(reg)
	if compile == nil {
		return
	}
	res = compile.FindAllString(s, -1)
	return
}

func FindAllString(s string, reg string) (res []string, err error) {
	compile, err := regexp.Compile(reg)
	if err != nil {
		return nil, err
	}
	res = compile.FindAllString(s, -1)
	return
}

// MatchNumber 提取字符串中的数字
func MatchNumber(s string) (res []string, err error) {
	compile, err := regexp.Compile("[0-9]+")
	if err != nil {
		return
	}
	if !compile.MatchString(s) {
		return
	}
	res = compile.FindAllString(s, -1)
	return
}

// MatchAllParams 匹配正则表达式中的参数值，args的参数个数必须与传入的正则表达式公式个数相等
func MatchAllParams(s string, reg string, args ...interface{}) (res []interface{}, err error) {
	compile, err := regexp.Compile(reg)
	if err != nil {
		return
	}
	if !compile.MatchString(s) {
		return
	}
	if match := compile.FindStringSubmatch(s); len(args) == len(match)-1 {
		for i := 0; i < len(args); i++ {
			switch args[i].(type) {
			case string:
				res = append(res, match[i+1])
			case int:
				parseInt, e := strconv.ParseInt(match[i+1], 10, 64)
				if e != nil {
					return
				}
				res = append(res, int(parseInt))
			case int64:
				parseInt, e := strconv.ParseInt(match[i+1], 10, 64)
				if e != nil {
					return
				}
				res = append(res, parseInt)
			case int32:
				parseInt, e := strconv.ParseInt(match[i+1], 10, 64)
				if e != nil {
					return
				}
				res = append(res, int32(parseInt))
			case bool:
				if match[i+1] == "true" {
					res = append(res, true)
				} else {
					res = append(res, false)
				}
			default:
				err = errors.New("no find type")
			}

		}
		return
	}
	err = errors.New("the number of parameters does not match")
	return
}
