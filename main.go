package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"shuZhiNet/model/student"
	"shuZhiNet/service/auth"
	"shuZhiNet/service/crawl"
	"shuZhiNet/service/login"
)

func allowOrigin(w http.ResponseWriter) http.ResponseWriter {
	w.Header().Set("Access-Control-Allow-Origin", "*")             //允许访问所有域
	w.Header().Add("Access-Control-Allow-Headers", "Content-Type") //header的类型
	w.Header().Set("content-type", "application/json")             //返回数据格式是json
	return w
}

// 登录 存贮用户登陆态 返回给前端姓名和JWT
// TODO: 前端姓名获取并不一定成功
func loginHandler(w http.ResponseWriter, r *http.Request) {
	w = allowOrigin(w)
	var input struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}
	body, _ := ioutil.ReadAll(r.Body)
	json.Unmarshal(body, &input)
	loginStudent := login.Login(input.Username, input.Password)
	var output struct {
		StudentName string `json:"student_name"`
		Token       string `json:"token"`
	}
	output.StudentName = loginStudent.Name
	output.Token = auth.GenerateJWT(loginStudent)
	student.Save(loginStudent)
	fmt.Println(output.StudentName, "logged in.")
	outputJSON, _ := json.Marshal(output)
	w.Write(outputJSON)
	fmt.Println(crawl.FetchActivitiesByStudent())
}

func engageHandler(w http.ResponseWriter, r *http.Request) {
	w = allowOrigin(w)
	var input struct {
		ID   string `json:"id"`
		hdid string `json:"hdid"`
		//nmd wsm 手机号码要取这个名，沿用了
		shouJhm string `json:"shouJhm"`
		email   string `json:"email"`
		Token   string `json:"token"`
	}
	body, _ := ioutil.ReadAll(r.Body)
	json.Unmarshal(body, &input)
	Loggedin := auth.AuthToken(input.ID, input.Token)
	if Loggedin {
		student, _ := student.Get(input.ID)
		var cookies []*http.Cookie
		cookies = append(cookies, &student.Cookie)
		//shuZhiNetUrl, _ := url.Parse("http://www.sz.shu.edu.cn")
		var jar http.CookieJar
		//jar.SetCookies(shuZhiNetUrl, cookies)
		client := http.Client{Jar: jar}
		fmt.Println(input.ID)
		fmt.Println(input.shouJhm)
		fmt.Println(input.Token)
		string1 := "http://www.sz.shu.edu.cn/api/HuoDong/HuoDBMXX/GetHuoDBM?hdid=" + input.hdid +
			"&shouJhm=" + input.shouJhm + "&email=" + input.email
		fmt.Println(string1)
		response, _ := client.Get(string1)
		res, _ := ioutil.ReadAll(response.Body)
		fmt.Println(res)
	}

}

func activitiesHandler(w http.ResponseWriter, r *http.Request) {
	w = allowOrigin(w)
	response, _ := json.Marshal(crawl.FetchActivitiesByStudent())
	w.Write(response)
}

func main() {
	http.HandleFunc("/login", loginHandler)
	http.HandleFunc("/activities", activitiesHandler)
	http.HandleFunc("/engage", engageHandler)
	err := http.ListenAndServe(":8000", nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
