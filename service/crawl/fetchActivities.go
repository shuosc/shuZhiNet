package crawl

import (
	"encoding/json"
	"github.com/PuerkitoBio/goquery"
	"io/ioutil"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"regexp"
	"shuZhiNet/infrastructure"
	"shuZhiNet/model/activity"
	"shuZhiNet/model/student"
	"strconv"
	"strings"
	"time"
)

// 得到最新的所有活动
/* HDType对应的活动分类
1--学业辅导
2--师生互动
3--志愿服务
4--社会实践
5--创新创业
6--文体活动
7--素质拓展
8--就业实习
*/
func FetchActivitiesByStudent() []activity.Activity {
	client := http.Client{}
	response, _ := client.Get("http://www.sz.shu.edu.cn/")
	content, _ := ioutil.ReadAll(response.Body)
	var activityList []activity.Activity
	activityType := 0
	reader := strings.NewReader(string(content))
	doc, _ := goquery.NewDocumentFromReader(reader)
	doc.Find(".ActiveScroll").Each(
		func(_ int, selection *goquery.Selection) {
			activityType++
			selection.Find(".ActiveScroll_bdli").Each(
				func(_ int, selection *goquery.Selection) {
					href, _ := selection.Find(".ActiveScroll_bdli > a").Attr("href")
					reHDid := regexp.MustCompile(`=([0-9]*)`)
					hdid := reHDid.Find([]byte(href))[1:]
					info := selection.Find(".ActiveInfo")
					// go's time format string is totally insane
					// faq
					startTime, _ := time.Parse("2006-01-02  15:04", info.Find(".icon_timeb").Text()[:17])
					endTime, _ := time.Parse("2006-01-02  15:04", info.Find(".icon_timeb").Text()[18:])
					signUpTime, _ := time.Parse("2006-01-02  15:04", info.Find(".TimeEnd").Text()[15:])
					thisActivity := activity.Activity{
						TypeId:     strconv.Itoa(activityType),
						Id:         string(hdid),
						Title:      info.Find(".ActiveTitle").Text(),
						Leader:     info.Find(".icon_leader").Text(),
						Address:    info.Find(".icon_addres").Text(),
						StartTime:  startTime,
						EndTime:    endTime,
						SignUpTime: signUpTime,
					}
					serialized, _ := json.Marshal(thisActivity)
					infrastructure.Redis.Set("Activity:"+string(hdid), serialized, 0)
					activity.Save(thisActivity)
					activityList = append(activityList, thisActivity)
				})
		})
	return activityList
}

type ParticipatedActivity struct {
	activity.Activity
	ParticipateInfoId string `json:"participate_info_id"`
}

func FetchParticipatingActivityIds(student student.Student) []ParticipatedActivity {
	jar, _ := cookiejar.New(nil)
	// todo: 处理多页
	participatingURL := "http://www.sz.shu.edu.cn/api/HuoDong/HuoDXX/GetHuoDBMXX?pageSize=30&pageNumber=1"
	participatingURLObject, _ := url.Parse(participatingURL)
	jar.SetCookies(participatingURLObject, student.Cookies)
	client := http.Client{Jar: jar}
	response, _ := client.Get(participatingURL)
	responseBody, _ := ioutil.ReadAll(response.Body)
	var responseJson struct {
		Data struct {
			ActivityInfo []struct {
				ParticipateInfoId int `json:"HuoDBMXXId"`
				Id                int `json:"Id"`
			} `json:"huodxx"`
		} `json:"data"`
	}
	json.Unmarshal(responseBody, &responseJson)
	var result []ParticipatedActivity
	for _, activityJson := range responseJson.Data.ActivityInfo {
		activityObject, _ := activity.Get(strconv.Itoa(activityJson.Id))
		result = append(result, ParticipatedActivity{
			Activity:          activityObject,
			ParticipateInfoId: strconv.Itoa(activityJson.ParticipateInfoId),
		})
	}
	return result
}
