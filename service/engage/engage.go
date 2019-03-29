package engage

import (
	"fmt"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"shuZhiNet/model/student"
)

func Engage(student student.Student, activityId string, phoneNumber string, mailAddress string) {
	jar, _ := cookiejar.New(nil)
	shuZhiNetUrl, _ := url.Parse("http://www.sz.shu.edu.cn")
	fmt.Println(student.Cookies)
	fmt.Println(shuZhiNetUrl)
	jar.SetCookies(shuZhiNetUrl, student.Cookies)
	client := http.Client{Jar: jar}
	string1 := "http://www.sz.shu.edu.cn/api/HuoDong/HuoDBMXX/GetHuoDBM?hdid=" + activityId +
		"&shouJhm=" + phoneNumber + "&email=" + mailAddress
	// let's use Get to post data to the server!
	// why not?
	// f**k shuZhiNet
	client.Get(string1)
}
