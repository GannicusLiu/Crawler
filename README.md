# Crawler
## Sina Weibo Login By Golang

# ����΢��ģ���¼

### ��¼��ַ
    http://weibo.com/login.php
    �Ѹ�ҳ���cookieȡ�����������¼�������ʱ����Ҫ�õ�
    
### ��ȡǰ�õ�¼�������
#### �����ַ
    http://login.sina.com.cn/sso/prelogin.php?entry=weibo&callback=sinaSSOController.preloginCallBack&su=&rsakt=mod&client=ssologin.js(v1.4.18)&_=1463016085877
#### ���ؽ��
    sinaSSOController.preloginCallBack({"retcode":0,"servertime":1463016089,"pcid":"gz-12a8a1ae997f9fafb79edd1445fae679f4a1","nonce":"JHVO2Z","pubkey":"EB2A38568661887FA180BDDB5CABD5F21C7BFD59C090CB2D245A87AC253062882729293E5506350508E7F9AA3BB77F4333231490F915F6D63C55FE2F08A49B353F444AD3993CACC02DB784ABBB8E42A9B1BBFFFB38BE18D78E87A0E41B9B8F73A928EE0CCEE1F6739884B9777E4FE9E88A1BBE495927AC4A799B3181D6442443","rsakv":"1330428213","uid":"1893131633","exectime":174})
    
#### ƥ�����
    ��������ʽ��ȡjson��������servertime��pcid��nonce��pubkey��rsakv
    - servertime ������ʱ��
    - pcid ���ص�¼��֤��ʱ���õ�
    - nonce ����� RSA��������ʱ���õ�
    - pubkey RSA���ܵĹ�Կ
    - rsakv ��¼����
    
### ������¼������װ
    http://login.sina.com.cn/js/sso/ssologin.js
#### ΢���û�����su�����ܹ���
    ��http://login.sina.com.cn/js/sso/ssologin.js 312��
    sinaSSOEncoder.base64.encode(urlencode(username));
#### ΢���û����루sp�����ܹ���
    ��http://login.sina.com.cn/js/sso/ssologin.js 902��
    if ((me.loginType & rsa) && me.servertime && sinaSSOEncoder && sinaSSOEncoder.RSAKey) {
    	request.servertime = me.servertime;
    	request.nonce = me.nonce;
    	request.pwencode = "rsa2";
    	request.rsakv = me.rsakv;
    	var RSAKey = new sinaSSOEncoder.RSAKey();
    	RSAKey.setPublic(me.rsaPubkey, "10001");
    	password = RSAKey.encrypt([me.servertime, me.nonce].join("\t") + "\n" + password)
    }
    
### ��֤���ַ
    http://login.sina.com.cn/cgi/pin.php?r=14749233&s=0&p=hk-4a80803307d47c997fa630d4109531dafc8e
    ����ǰ�õ�¼��ȡ���Ĳ���showpin
    ��showpinΪ1ʱ��Ҫ��֤�룬��Ϊ0ʱ����Ҫ��֤��
    ��Ҫ��֤��ʱ������֤��ͼ�����������ֹ���������Զ�ʶ��   
    
### ��¼
#### �����ַ
    http://login.sina.com.cn/sso/login.php?client=ssologin.js(v1.4.18)
#### �������
    entry=weibo
    gateway=1
    from=
    savestate=0
    useticket=1
    pagerefer=
    vsnf=1
    su=��������΢���û���
    service=miniblog
    servertime=ǰ�õ�¼����
    nonce=ǰ�õ�¼����
    pwencode=rsa2
    rsakv=ǰ�õ�¼����
    sp=���ܺ������
    sr=1366*768
    encoding=UTF-8
    prelt=1279
    url=http%3A%2F%2Fweibo.com%2Fajaxlogin.php%3Fframelogin%3D1%26callback%3Dparent.sinaSSOController.feedBackUrlCallBack
    returntype=META
    door=��֤�� ��Ҫ��֤���ʱ��Ŵ��ݸò���
    
#### ���ؽ��
    `location.replace('http://passport.weibo.com/wbsso/login?url=http%3A%2F%2Fweibo.com%2Fajaxlogin.php%3Fframelogin%3D1%26callback%3Dparent.sinaSSOController.feedBackUrlCallBack%26sudaref%3Dweibo.com&ticket=ST-MTg5MzEzMTYzMw==-1464838377-gz-E4356F083C2998D105269B1FA6BE6F5A&retcode=0');});}`
    ��location.replace�е�urlƥ�����
    
### ����passport
#### ����url
    ��¼���صĽ����ƥ�������url
#### ���ؽ��
    �����󷵻ص�header���и�302��ת
    Location: http://weibo.com/ajaxlogin.php?framelogin=1&callback=parent.sinaSSOController.feedBackUrlCallBack&sudaref=weibo.com

### 302��ת
#### ��������
    ��������������ݣ���ʾ��¼�ɹ��� uniqueid����������΢����ȫ��ΨһID,��������ȥ��
    <html><head><script language='javascript'>parent.sinaSSOController.feedBackUrlCallBack({"result":true,"userinfo":{"uniqueid":"xxx","userid":null,"displayname":null,"userdomain":"?wvr=5&lf=reg"}});</script></head><body></body></html>
    
### Ȼ����Ϳ���ȥ΢����������������
    ps���ǵð�http://login.sina.com.cn/sso/login.php?client=ssologin.js(v1.4.18)��post���󷵻ص�cookie����������ץȡ����΢��ȫ����Щcookie��
    
