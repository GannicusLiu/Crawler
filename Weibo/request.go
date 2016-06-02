/*
 * @functional golang http request
 * @author junchen168@live.cn
 */
package main

import (
	"crypto/tls"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"regexp"
	"strings"
)

/*
 * @functional 发送http请求
 * @param string method 请求的方法类型
 * @param string reqUrl 请求的url
 * @param string params 请求的参数
 * @param string cookies cookie字符串
 * @param string domain cookie的域
 * @param map mapHeader 请求的头信息
 * @return string
 */
func DoRequest(method, reqUrl, params, cookies, domain string, mapHeader map[string]string) (respBody, respCookies string) {
	return organizeRequest(method, reqUrl, params, cookies, domain, mapHeader, true, tls.VersionTLS10, tls.VersionTLS12)
}

/*
 * @functional 构造http 请求
 * @param string method 请求的方法类型
 * @param string reqUrl 请求的url
 * @param string params 参数
 * @param string cookies cookie字符串
 * @param string domain cookie的域
 * @param map mapHeader 请求的头信息
 * @param bool 是否启用tls
 * @param uint16 min tls
 * @param uint16 max tls
 * @return string
 */
func organizeRequest(method, reqUrl, params, cookies, domain string, mapHeader map[string]string, tlsFlag bool, mintls, maxtls uint16) (respBody, respCookies string) {
	method = strings.ToUpper(method)
	var strParams io.Reader
	if method == "POST" {
		strParams = strings.NewReader(params)
	} else {
		strParams = nil
	}

	req, err := http.NewRequest(method, reqUrl, strParams)
	if err != nil {
		log.Fatalln("NewRequest Err:", err)
		return
	}

	if mapHeader != nil {
		for key, val := range mapHeader {
			req.Header.Set(key, val)
		}
	}

	gCookieJar, _ := cookiejar.New(nil)
	if len(cookies) > 0 {
		cookies := appendCookies(cookies, "", domain)
		cookieUrl, _ := url.Parse(reqUrl)
		gCookieJar.SetCookies(cookieUrl, cookies)
	}
	var transport *http.Transport = nil
	proxyUrl, _ := url.Parse("http://127.0.0.1:8888")
	if tlsFlag {
		transport = &http.Transport{
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: true,
				MinVersion:         mintls,
				MaxVersion:         maxtls,
			},
			DisableCompression:true,
			Proxy: http.ProxyURL(proxyUrl),
		}
	} else {
		transport = &http.Transport{}
	}

	client := &http.Client{
		Transport:     transport,
		CheckRedirect: nil,
		Jar:           gCookieJar,
	}

	resp, err := client.Do(req)
	if err != nil {
		log.Fatalln("Do Request Err:", err)
		return
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatalln("ReadAll Response Err:", err)
		return
	}
	respBody = string(body)

	arrCookies := resp.Cookies()
	for _, data := range arrCookies {
		respCookies += data.Name + "=" + data.Value + ";"
	}
	if len(respCookies) > 0 {
		respCookies = SubString(respCookies, 0, len(respCookies)-1)
	}

	return
}

/*
 * @functional http请求附加cookie
 * @param string strCookies
 * @return []*http.Cookie
 */
func appendCookies(strCookies, path, domain string) []*http.Cookie {
	var cookies []*http.Cookie

	if path == "" {
		path = "/"
	}

	mapCookie := saveCookies(strCookies)
	for k, v := range mapCookie {
		appendCookie := &http.Cookie{
			Name:   k,
			Value:  v,
			Path:   path,
			Domain: domain,
		}
		cookies = append(cookies, appendCookie)
	}
	return cookies
}

/*
 * @functional 解析cookie到map
 * @param string cookies cookie字符串
 * @return map
 */
func saveCookies(cookies string) map[string]string {
	mapCookie := make(map[string]string)
	reg := regexp.MustCompile(`([^=]+)=([^;]*);?`)
	arrCookie := reg.FindAllStringSubmatch(cookies, -1)
	if len(arrCookie) > 0 {
		for i := 0; i < len(arrCookie); i++ {
			mapCookie[arrCookie[i][1]] = arrCookie[i][2]
		}
	}
	return mapCookie
}
