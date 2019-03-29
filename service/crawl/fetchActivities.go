package crawl

import (
	"github.com/PuerkitoBio/goquery"
	"io/ioutil"
	"net/http"
	"regexp"
	"shuZhiNet/model/activity"
	"strings"
	"time"
)

// 得到最新的所有活动
func FetchActivitiesByStudent() []activity.Activity {
	client := http.Client{}
	response, _ := client.Get("http://www.sz.shu.edu.cn/")
	content, _ := ioutil.ReadAll(response.Body)
	var HDList []activity.Activity
	reader := strings.NewReader(string(content))
	doc, _ := goquery.NewDocumentFromReader(reader)
	doc.Find(".ActiveScroll_bdli").Each(
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
				Id:         string(hdid),
				Title:      info.Find(".ActiveTitle").Text(),
				Leader:     info.Find(".icon_leader").Text(),
				Address:    info.Find(".icon_addres").Text(),
				StartTime:  startTime,
				EndTime:    endTime,
				SignUpTime: signUpTime,
			})
		})
	return HDList
}
