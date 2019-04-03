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
	fmt.Println(student)
	shuZhiNetUrl, _ := url.Parse("http://www.sz.shu.edu.cn")
	jar.SetCookies(shuZhiNetUrl, student.Cookies)
	client := http.Client{Jar: jar}
	getURl := "http://www.sz.shu.edu.cn/api/HuoDong/HuoDBMXX/GetHuoDBM?hdid=" + activityId +
		"&shouJhm=" + phoneNumber + "&email=" + mailAddress
	client.Get(getURl)
	fmt.Print(student.Id, " engaged in ", activityId)
}
