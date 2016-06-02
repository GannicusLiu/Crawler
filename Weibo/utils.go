/*
 * @functional 工具函数
 * @author junchen168@live.cn
 */
package main

import (
	"regexp"

	"os"

	"io/ioutil"
)

/*
 * @functional 正则表达式提取数据
 * @param string strText 输入文本
 * @param string strReg 正则表达式
 * @return string
 */
func RegexFind(strText, strReg string) (result string) {
	reg := regexp.MustCompile(strReg)
	arrMatch := reg.FindAllStringSubmatch(strText, -1)
	if len(arrMatch) > 0 {
		result = arrMatch[0][1]
	}
	return
}

/*
 * @functional截取字符串
 * @param string str 原始字符串
 * @param int begin 截取开始位置
 * @param int length 截取长度
 * @return string
 */
func SubString(str string, begin, length int) (substr string) {
	rs := []rune(str)
	lth := len(rs)
	if begin < 0 {
		begin = 0
	}
	if begin >= lth {
		begin = lth
	}
	end := begin + length
	if end > lth {
		end = lth
	}

	return string(rs[begin:end])
}

/*
 * @functional判断文件夹是否存在
 * @param string path 文件路径
 * @return bool
 */
func IsDirExist(path string) bool {
	p, err := os.Stat(path)
	if err != nil {
		return os.IsExist(err)
	} else {
		return p.IsDir()
	}
}

/*
 * @functional 读取文件
 * @param string path 文件路径
 * @return bool
 */
func ReadFile(path string) string {
	fi, err := os.Open(path)
	if err != nil {
		return ""
	}
	defer fi.Close()
	fd, err := ioutil.ReadAll(fi)
	return string(fd)
}

/**
 * @functional 正则表达式匹配数据
 * @string strText 源字符串
 * @string strReg 正则表达式
 * @return string
 */
func MatchData(strText, strReg string) (result string) {
	reg := regexp.MustCompile(strReg)
	arrMatch := reg.FindAllStringSubmatch(strText, -1)
	if len(arrMatch) > 0 {
		result = arrMatch[0][1]
	}
	return
}
