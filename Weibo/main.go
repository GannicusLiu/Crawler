/*
 * @functional 新浪微博登录
 * @author junchen168@live.cn
 */
package main

import (
	"crypto/rand"
	"crypto/rsa"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"math/big"
	"net/url"
	"strconv"

	"bytes"
	"io"
	"os"
	"time"

	"github.com/bitly/go-simplejson"
)

var (
	uname       = "微博用户名"
	password    = "微博密码"
	captchaPath = "../captcha/"
)

//编码用户名
func encryptUname(uname string) string {
	urlEncode := url.QueryEscape(uname)
	return base64.StdEncoding.EncodeToString([]byte(urlEncode))
}

//把字符串转换bigint
func string2big(s string) *big.Int {
	ret := new(big.Int)
	ret.SetString(s, 16)
	return ret
}

//加密密码
func encryptPassword(loginInfo map[string]string, password string) string {
	pub := rsa.PublicKey{
		N: string2big(loginInfo["pubkey"]),
		E: 65537,
	}
	encryString := loginInfo["servertime"] + "\t" + loginInfo["nonce"] + "\n" + password
	encryResult, _ := rsa.EncryptPKCS1v15(rand.Reader, &pub, []byte(encryString))
	return hex.EncodeToString(encryResult)
}

//获取登录页面的cookie
func getLoginPageCookies() (strCookies string) {
	strLoginUrl := `http://weibo.com/login.php`
	_, strCookies = DoRequest(`GET`, strLoginUrl, ``, ``, ``, nil)
	return
}

//获取登录参数
func getLoginInfo(su string) (strJson, cookies string) {
	preUrl := `http://login.sina.com.cn/sso/prelogin.php?entry=weibo&callback=sinaSSOController.preloginCallBack&su=` + su + `&rsakt=mod&checkpin=1&client=ssologin.js(v1.4.18)&_=`

	resp, cookies := DoRequest(`GET`, preUrl, ``, ``, ``, nil)
	strJson = RegexFind(resp, `\((.*?)\)`)
	return
}

//解析json
func loadJson(strJson string) (loginInfo map[string]string) {
	jsonParser, _ := simplejson.NewJson([]byte(strJson))
	loginInfo = make(map[string]string)
	unixtime, _ := jsonParser.Get("servertime").Int()
	loginInfo["servertime"] = strconv.Itoa(unixtime)
	loginInfo["rsakv"], _ = jsonParser.Get("rsakv").String()
	loginInfo["pcid"], _ = jsonParser.Get("pcid").String()
	loginInfo["nonce"], _ = jsonParser.Get("nonce").String()
	loginInfo["pubkey"], _ = jsonParser.Get("pubkey").String()
	showpin, _ := jsonParser.Get("showpin").Int()
	loginInfo["showpin"] = strconv.Itoa(showpin)
	return
}

//保存验证码
func saveCaptcha(pcid, cookies string) {
	rnd := time.Now().Format("20060102150405")
	captchUrl := "http://login.sina.com.cn/cgi/pin.php?r=" + rnd + "&s=0&p=" + pcid
	captcha, _ := DoRequest(`GET`, captchUrl, ``, cookies, ``, nil)
	imgSave, err := os.Create(captchaPath + rnd + ".png")
	if err != nil {
		fmt.Println(err.Error())
	}
	io.Copy(imgSave, bytes.NewReader([]byte(captcha)))
}

//开始登录
func weiboLogin(su, sp, captcha, cookies string, loginInfo map[string]string) (postResp, postCookies string) {
	strPostUrl := `http://login.sina.com.cn/sso/login.php?client=ssologin.js(v1.4.18)`
	var strParams = "entry=weibo&gateway=1&from=&savestate=0&useticket=1&pagerefer=&vsnf=1&su=" + su + "&service=miniblog&servertime=" + loginInfo["servertime"] + "&nonce=" + loginInfo["nonce"] + "&pwencode=rsa2&rsakv=" + loginInfo["rsakv"] + "&sp=" + sp + "&sr=1366*768&encoding=UTF-8&prelt=1279&url=http%3A%2F%2Fweibo.com%2Fajaxlogin.php%3Fframelogin%3D1%26callback%3Dparent.sinaSSOController.feedBackUrlCallBack&returntype=META"
	//需要验证码
	if loginInfo["showpin"] == "1" {
		strParams += "&door=" + captcha
	}
	header := map[string]string{
		"Host":                      "login.sina.com.cn",
		"Proxy-Connection":          "keep-alive",
		"Cache-Control":             "max-age=0",
		"Accept":                    "text/html,application/xhtml+xml,application/xml;q=0.9,image/webp,*/*;q=0.8",
		"Origin":                    "http://weibo.com",
		"Upgrade-Insecure-Requests": "1",
		"User-Agent":                "Mozilla/5.0 (Windows NT 10.0; WOW64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/50.0.2661.94 Safari/537.36",
		"Referer":                   "http://weibo.com",
		"Accept-Language":           "zh-CN,zh;q=0.8,en;q=0.6,ja;q=0.4",
		"Content-Type":              "application/x-www-form-urlencoded",
	}
	postResp, postCookies = DoRequest(`POST`, strPostUrl, strParams, cookies, ``, header)
	return
}

//获取passport并请求
func callPassport(resp, cookies string) (passresp, passcookies string) {
	//提取passport跳转地址
	passportUrl := RegexFind(resp, `location.replace\('(.*?)'\)`)
	header := map[string]string{
		"Host":                      "passport.weibo.com",
		"Proxy-Connection":          "keep-alive",
		"Cache-Control":             "max-age=0",
		"Accept":                    "text/html,application/xhtml+xml,application/xml;q=0.9,image/webp,*/*;q=0.8",
		"Origin":                    "http://weibo.com",
		"Upgrade-Insecure-Requests": "1",
		"User-Agent":                "Mozilla/5.0 (Windows NT 10.0; WOW64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/50.0.2661.94 Safari/537.36",
		"Referer":                   "http://login.sina.com.cn/sso/login.php?client=ssologin.js(v1.4.18)",
	}
	passresp, passcookies = DoRequest(`GET`, passportUrl, ``, cookies, ``, header)
	return
}

//进入首页
func entryHome(redirectUrl, cookies string) (homeResp, homeCookies string) {
	header := map[string]string{
		"Host":                      "weibo.com",
		"Connection":                "keep-alive",
		"Accept":                    "text/html,application/xhtml+xml,application/xml;q=0.9,image/webp,*/*;q=0.8",
		"Origin":                    "http://weibo.com",
		"User-Agent":                "Mozilla/5.0 (Windows NT 10.0; WOW64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/50.0.2661.94 Safari/537.36",
		"Upgrade-Insecure-Requests": "1",
		"Referer":                   "http://weibo.com/",
	}
	homeResp, homeCookies = DoRequest(`GET`, redirectUrl, ``, cookies, ``, header)
	return
}

//抓取我的微博页面
func getPage(cookies string) (resp string) {
	url := "http://weibo.com/onfoucs/profile?rightmod=1&wvr=6&mod=personinfo&is_all=1"
	header := map[string]string{
		"Host":                      "weibo.com",
		"Connection":                "keep-alive",
		"Accept":                    "text/html,application/xhtml+xml,application/xml;q=0.9,image/webp,*/*;q=0.8",
		"Origin":                    "http://weibo.com",
		"User-Agent":                "Mozilla/5.0 (Windows NT 10.0; WOW64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/50.0.2661.94 Safari/537.36",
		"Upgrade-Insecure-Requests": "1",
		"Referer":                   "http://weibo.com/",
	}
	resp, _ = DoRequest(`GET`, url, ``, cookies, ``, header)
	return
}

func main() {
	if !IsDirExist(captchaPath) {
		os.Mkdir(captchaPath, 0755)
	}
	//base64微博用户名
	su := encryptUname(uname)
	//获取登录页面cookies
	cookies := getLoginPageCookies()
	//提取登录用户到的servertime、nonce、rsakv、showpin
	strJson, _ := getLoginInfo(su)
	loginInfo := loadJson(strJson)
	/*
		showpin为1时表示需要输入验证码，此处是下载验证码图片并阻塞程序直到手动输入验证码;
		在生产环境中可以直接将图片流传到客户端让用户输入或者调用打码平台自动识别
	*/
	var captcha string
	if loginInfo["showpin"] == "1" {
		saveCaptcha(loginInfo["pcid"], cookies)
		inputDone := make(chan string)
		go func() {
			for {
				fmt.Println("waiting for input captcha...")
				input := ReadFile(captchaPath + "captcha.txt")
				if input != "" {
					inputDone <- input
					break
				}
			}
		}()
		captcha = <-inputDone
	}

	sp := encryptPassword(loginInfo, password)
	postResp, loginCookies := weiboLogin(su, sp, captcha, cookies, loginInfo)
	//请求passport
	passportResp, _ := callPassport(postResp, cookies+";"+loginCookies)
	uniqueid := MatchData(passportResp, `"uniqueid":"(.*?)"`)
	homeUrl := "http://weibo.com/u/" + uniqueid + "/home?topnav=1&wvr=6"
	//进入个人主页
	entryHome(homeUrl, loginCookies)
	//抓取个首页
	result := getPage(loginCookies)
	fmt.Println(result)
}
