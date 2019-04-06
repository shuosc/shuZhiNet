package student

import (
	"github.com/PuerkitoBio/goquery"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"shuZhiNet/model/student"
)

func getStudentName(jar http.CookieJar) string {
	client := http.Client{Jar: jar}
	nameResponse, _ := client.Get("http://www.sz.shu.edu.cn/index.aspx")
	nameDoc, _ := goquery.NewDocumentFromReader(nameResponse.Body)
	name := nameDoc.Find("#UserLoginUserControl_lbUserName").Text()
	return name
}

func Login(username string, password string) (student.Student, error) {
	jar, _ := cookiejar.New(nil)
	client := http.Client{Jar: jar}
	authData := url.Values{}
	authData.Add("password", password)
	authData.Add("userName", username)
	_, err := client.PostForm("http://www.sz.shu.edu.cn/api/Sys/Users/Login", authData)
	if err != nil {
		return student.Student{}, nil
	}
	shuZhiNetUrl, _ := url.Parse("http://www.sz.shu.edu.cn")
	return student.Student{
		Id:      username,
		Name:    getStudentName(client.Jar),
		Cookies: jar.Cookies(shuZhiNetUrl),
	}, nil
}
