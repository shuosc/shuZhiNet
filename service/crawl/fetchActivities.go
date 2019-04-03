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
	var HDList []activity.Activity
	HDType := 0
	reader := strings.NewReader(string(content))
	doc, _ := goquery.NewDocumentFromReader(reader)
	doc.Find(".ActiveScroll").Each(
		func(_ int, selection *goquery.Selection) {
			HDType = HDType + 1
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
						TypeId:     strconv.Itoa(HDType),
						Id:         string(hdid),
						Title:      info.Find(".ActiveTitle").Text(),
						Leader:     info.Find(".icon_leader").Text(),
						Address:    info.Find(".icon_addres").Text(),
						StartTime:  startTime,
						EndTime:    endTime,
						SignUpTime: signUpTime,
					}
					serilized, _ := json.Marshal(thisActivity)
					infrastructure.Redis.Set("Activity:"+string(hdid), serilized, 0)
					HDList = append(HDList, thisActivity)
				})
		})
	return HDList
}

//得到学生已经参与的活动
func FetchMyActivity(student student.Student) ([]string, []string) {
	var cancelIdList []string
	jar, _ := cookiejar.New(nil)
	cancelListURL := "http://www.sz.shu.edu.cn/api/HuoDong/HuoDXX/GetHuoDBMXX?pageSize=3&pageNumber=1" /**/
	MyActivity, _ := url.Parse(cancelListURL)
	jar.SetCookies(MyActivity, student.Cookies)
	client := http.Client{Jar: jar}
	response, _ := client.Get(cancelListURL)
	content, _ := ioutil.ReadAll(response.Body)
	reader := strings.NewReader(string(content))
	doc, _ := goquery.NewDocumentFromReader(reader)
	reDBid := regexp.MustCompile(`XXId":([0-9]*)`)
	DBid := reDBid.FindAllString(doc.Text(), -1)
	reActivityId := regexp.MustCompile(`"Id":([0-9]*)`)
	activityId := reActivityId.FindAllString(doc.Text(), -1)
	for i := 0; i < len(DBid); i++ {
		activityId[i] = activityId[i][5:]
		cancelIdList = append(cancelIdList, DBid[i][6:])
	}
	return activityId, cancelIdList
}
