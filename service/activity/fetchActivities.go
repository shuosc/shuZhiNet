package activity

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
	"time"
)

func FetchActivities() []activity.Activity {
	response, _ := http.Get("http://www.sz.shu.edu.cn/")
	var activityList []activity.Activity
	doc, _ := goquery.NewDocumentFromReader(response.Body)
	typeNames := doc.Find(".MainTab_hdli.Tabli>em").Map(func(_ int, selection *goquery.Selection) string {
		return selection.Text()
	})
	var types []activity.Type
	for _, typeName := range typeNames {
		types = append(types, activity.GetTypeByName(typeName))
	}
	doc.Find(".MainTab_bdul.h417.Tabul .ActiveScroll_bd ul").Each(func(index int, outerSelection *goquery.Selection) {
		typeObject := types[index]
		outerSelection.Find(".ActiveScroll_bdli").Each(func(_ int, selection *goquery.Selection) {
			href, _ := selection.Find("a").Attr("href")
			activityIdRegex := regexp.MustCompile("hdid=(\\d+)")
			hdId := activityIdRegex.FindString(href)[5:]
			info := selection.Find(".ActiveInfo")
			startTime, _ := time.Parse("2006-01-02  15:04", info.Find(".icon_timeb").Text()[:17])
			endTime, _ := time.Parse("2006-01-02  15:04", info.Find(".icon_timeb").Text()[18:])
			signUpTime, _ := time.Parse("2006-01-02  15:04", info.Find(".TimeEnd").Text()[15:])
			thisActivity := activity.Activity{
				TypeId:     typeObject.Id,
				Id:         hdId,
				Title:      info.Find(".ActiveTitle").Text(),
				Leader:     info.Find(".icon_leader").Text(),
				Address:    info.Find(".icon_addres").Text(),
				StartTime:  startTime,
				EndTime:    endTime,
				SignUpTime: signUpTime,
			}
			serialized, _ := json.Marshal(thisActivity)
			infrastructure.Redis.Set("Activity:"+string(hdId), serialized, 0)
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

func FetchParticipatingActivities(student student.Student) []ParticipatedActivity {
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
				ParticipateInfoId int    `json:"HuoDBMXXId"`
				Id                int    `json:"Id"`
				Title             string `json:"HuoDMC"`
				Address           string `json:"HuoDDD"`
				StartTime         string `json:"HuoDKSSJ"`
				EndTime           string `json:"HuoDJSSJ"`
				SignUpTime        string `json:"CreatedOn"`
				TypeName          string `json:"HuoDLBName"`
				Leader            string `json:"ZuZMC"`
			} `json:"huodxx"`
		} `json:"data"`
	}
	_ = json.Unmarshal(responseBody, &responseJson)
	var result []ParticipatedActivity
	for _, activityJson := range responseJson.Data.ActivityInfo {
		activityObject, err := activity.Get(strconv.Itoa(activityJson.Id))
		if err != nil {
			startTime, _ := time.Parse("2006-01-02T15:04:05", activityJson.StartTime)
			endTime, _ := time.Parse("2006-01-02T15:04:05", activityJson.EndTime)
			signUpTime, _ := time.Parse("2006-01-02T15:04:05", activityJson.SignUpTime)
			activityObject = activity.Activity{
				Id:         strconv.Itoa(activityJson.Id),
				TypeId:     activity.GetTypeByName(activityJson.TypeName).Id,
				Title:      activityJson.Title,
				Leader:     activityJson.Leader,
				Address:    activityJson.Address,
				StartTime:  startTime,
				EndTime:    endTime,
				SignUpTime: signUpTime,
			}
			activity.Save(activityObject)
		}
		result = append(result, ParticipatedActivity{
			Activity:          activityObject,
			ParticipateInfoId: strconv.Itoa(activityJson.ParticipateInfoId),
		})
	}
	return result
}
