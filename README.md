# Crawler
## Sina Weibo Login By Golang

# 新浪微博模拟登录

### 登录地址
    http://weibo.com/login.php
    把该页面的cookie取下来，后面登录发请求的时候需要用到
    
### 获取前置登录所需参数
#### 请求地址
    http://login.sina.com.cn/sso/prelogin.php?entry=weibo&callback=sinaSSOController.preloginCallBack&su=&rsakt=mod&client=ssologin.js(v1.4.18)&_=1463016085877
#### 返回结果
    sinaSSOController.preloginCallBack({"retcode":0,"servertime":1463016089,"pcid":"gz-12a8a1ae997f9fafb79edd1445fae679f4a1","nonce":"JHVO2Z","pubkey":"EB2A38568661887FA180BDDB5CABD5F21C7BFD59C090CB2D245A87AC253062882729293E5506350508E7F9AA3BB77F4333231490F915F6D63C55FE2F08A49B353F444AD3993CACC02DB784ABBB8E42A9B1BBFFFB38BE18D78E87A0E41B9B8F73A928EE0CCEE1F6739884B9777E4FE9E88A1BBE495927AC4A799B3181D6442443","rsakv":"1330428213","uid":"1893131633","exectime":174})
    
#### 匹配参数
    用正则表达式获取json包解析出servertime、pcid、nonce、pubkey、rsakv
    - servertime 服务器时间
    - pcid 下载登录验证码时会用到
    - nonce 随机串 RSA加密密码时会用到
    - pubkey RSA加密的公钥
    - rsakv 登录参数
    
### 分析登录参数组装
    http://login.sina.com.cn/js/sso/ssologin.js
#### 微博用户名（su）加密规则
    见http://login.sina.com.cn/js/sso/ssologin.js 312行
    sinaSSOEncoder.base64.encode(urlencode(username));
#### 微博用户密码（sp）加密规则
    见http://login.sina.com.cn/js/sso/ssologin.js 902行
    if ((me.loginType & rsa) && me.servertime && sinaSSOEncoder && sinaSSOEncoder.RSAKey) {
    	request.servertime = me.servertime;
    	request.nonce = me.nonce;
    	request.pwencode = "rsa2";
    	request.rsakv = me.rsakv;
    	var RSAKey = new sinaSSOEncoder.RSAKey();
    	RSAKey.setPublic(me.rsaPubkey, "10001");
    	password = RSAKey.encrypt([me.servertime, me.nonce].join("\t") + "\n" + password)
    }
    
### 验证码地址
    http://login.sina.com.cn/cgi/pin.php?r=14749233&s=0&p=hk-4a80803307d47c997fa630d4109531dafc8e
    根据前置登录获取到的参数showpin
    当showpin为1时需要验证码，当为0时不需要验证码
    需要验证码时，把验证码图保存下来，手工输入或者自动识别   
    
### 登录
#### 请求地址
    http://login.sina.com.cn/sso/login.php?client=ssologin.js(v1.4.18)
#### 请求参数
    entry=weibo
    gateway=1
    from=
    savestate=0
    useticket=1
    pagerefer=
    vsnf=1
    su=编码过后的微博用户名
    service=miniblog
    servertime=前置登录参数
    nonce=前置登录参数
    pwencode=rsa2
    rsakv=前置登录参数
    sp=机密后的密码
    sr=1366*768
    encoding=UTF-8
    prelt=1279
    url=http%3A%2F%2Fweibo.com%2Fajaxlogin.php%3Fframelogin%3D1%26callback%3Dparent.sinaSSOController.feedBackUrlCallBack
    returntype=META
    door=验证码 需要验证码的时候才传递该参数
    
#### 返回结果
    `location.replace('http://passport.weibo.com/wbsso/login?url=http%3A%2F%2Fweibo.com%2Fajaxlogin.php%3Fframelogin%3D1%26callback%3Dparent.sinaSSOController.feedBackUrlCallBack%26sudaref%3Dweibo.com&ticket=ST-MTg5MzEzMTYzMw==-1464838377-gz-E4356F083C2998D105269B1FA6BE6F5A&retcode=0');});}`
    把location.replace中的url匹配出来
    
### 请求passport
#### 请求url
    登录返回的结果中匹配出来的url
#### 返回结果
    该请求返回的header中有个302跳转
    Location: http://weibo.com/ajaxlogin.php?framelogin=1&callback=parent.sinaSSOController.feedBackUrlCallBack&sudaref=weibo.com

### 302跳转
#### 返回内容
    如果返回如下内容，表示登录成功了 uniqueid是我在新浪微博的全局唯一ID,我这里隐去了
    <html><head><script language='javascript'>parent.sinaSSOController.feedBackUrlCallBack({"result":true,"userinfo":{"uniqueid":"xxx","userid":null,"displayname":null,"userdomain":"?wvr=5&lf=reg"}});</script></head><body></body></html>
    
### 然后你就可以去微博做爱做的事情了
    ps：记得把http://login.sina.com.cn/sso/login.php?client=ssologin.js(v1.4.18)该post请求返回的cookie保存下来，抓取新浪微博全靠这些cookie了
    
