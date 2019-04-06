package activity

import (
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"shuZhiNet/model/student"
)

func TakePartInActivity(student student.Student, activityId string, phoneNumber string, mailAddress string) {
	jar, _ := cookiejar.New(nil)
	shuZhiNetUrl, _ := url.Parse("http://www.sz.shu.edu.cn")
	jar.SetCookies(shuZhiNetUrl, student.Cookies)
	client := http.Client{Jar: jar}
	getURl := "http://www.sz.shu.edu.cn/api/HuoDong/HuoDBMXX/GetHuoDBM?hdid=" + activityId +
		"&shouJhm=" + phoneNumber + "&email=" + mailAddress
	_, _ = client.Get(getURl)
}
