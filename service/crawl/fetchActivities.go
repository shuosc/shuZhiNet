package crawl

import (
	"github.com/PuerkitoBio/goquery"
	"io/ioutil"
	"net/http"
	"regexp"
	"shuZhiNet/model/activity"
	"strings"
	"time"
	"strconv"
)

// 得到最新的所有活动
func FetchActivitiesByStudent() []activity.Activity {
	client := http.Client{}
	response, _ := client.Get("http://www.sz.shu.edu.cn/")
	content, _ := ioutil.ReadAll(response.Body)
	var HDList []activity.Activity
	HDType := 0
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
					HDList = append(HDList, activity.Activity{
						TypeId:     strconv.Itoa(HDType),
						Id:         string(hdid),
						Title:      info.Find(".ActiveTitle").Text(),
						Leader:     info.Find(".icon_leader").Text(),
						Address:    info.Find(".icon_addres").Text(),
						StartTime:  startTime,
						EndTime:    endTime,
						SignUpTime: signUpTime,
					})
				})
		})
	return HDList
}
