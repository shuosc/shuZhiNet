package scholarship

import (
	"github.com/PuerkitoBio/goquery"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"regexp"
	"shuZhiNet/model/student"
	"strconv"
	"strings"
)

var Types = map[int]string{
	0: "学业优秀奖学金",
	4: "领导力奖学金",
	5: "创新创业奖学金",
	6: "公益爱心奖学金",
	7: "文艺体育奖学金",
	8: "自强不息奖学金",
}

func postFormWithUA(client *http.Client, url string, form url.Values) (*http.Response, error) {
	const UA string = "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_14_3) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/73.0.3683.86 Safari/537.36"
	req, _ := http.NewRequest("POST", url, strings.NewReader(form.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("User-Agent", UA)
	return client.Do(req)
}

func Apply(student student.Student, scholarshipType int, qualifications []int, reason string) {
	client := student.Client()
	applyPageResponse, _ := client.Get("http://www.sz.shu.edu.cn/JiangXj/JiangXJStudentList.aspx")
	applyPageDoc, _ := goquery.NewDocumentFromReader(applyPageResponse.Body)
	eventValidation, _ := applyPageDoc.Find("input[name=__EVENTVALIDATION]").Attr("value")
	viewState, _ := applyPageDoc.Find("input[name=__VIEWSTATE]").Attr("value")
	viewStateGenerator, _ := applyPageDoc.Find("input[name=__VIEWSTATEGENERATOR]").Attr("value")
	form := url.Values{
		"ScriptManager1":       []string{"UpdatePanel1|JiangXJddl"},
		"__EVENTTARGET":        []string{"JiangXJddl"},
		"__EVENTARGUMENT":      []string{""},
		"__LASTFOCUS":          []string{""},
		"__VIEWSTATE":          []string{viewState},
		"__VIEWSTATEGENERATOR": []string{viewStateGenerator},
		"__EVENTVALIDATION":    []string{eventValidation},
		"JiangXJddl":           []string{strconv.Itoa(scholarshipType)},
		"LiYtxt":               []string{""},
		"__ASYNCPOST":          []string{"true"},
	}
	typePageResponse, _ := postFormWithUA(&client, "http://www.sz.shu.edu.cn/JiangXj/JiangXJStudentList.aspx", form)
	typePageBody, _ := ioutil.ReadAll(typePageResponse.Body)
	viewStateRegex, _ := regexp.Compile("__VIEWSTATE\\|([^|]+)\\|")
	viewState = string(viewStateRegex.FindSubmatch(typePageBody)[1])
	viewStateGeneratorRegex, _ := regexp.Compile("__VIEWSTATEGENERATOR\\|([^|]+)\\|")
	viewStateGenerator = string(viewStateGeneratorRegex.FindSubmatch(typePageBody)[1])
	eventValidationRegex, _ := regexp.Compile("__EVENTVALIDATION\\|([^|]+)\\|")
	eventValidation = string(eventValidationRegex.FindSubmatch(typePageBody)[1])
	form = url.Values{
		"ScriptManager1":       []string{"UpdatePanel1|Savebt"},
		"JiangXJddl":           []string{strconv.Itoa(scholarshipType)},
		"LiYtxt":               []string{reason},
		"__EVENTTARGET":        []string{""},
		"__EVENTARGUMENT":      []string{""},
		"__LASTFOCUS":          []string{""},
		"__VIEWSTATE":          []string{viewState},
		"__VIEWSTATEGENERATOR": []string{viewStateGenerator},
		"__EVENTVALIDATION":    []string{eventValidation},
		"__ASYNCPOST":          []string{"true"},
		"Savebt":               []string{"提   交"},
	}
	for _, qualification := range qualifications {
		qualificationIdRegex, _ := regexp.Compile("name=\"JiangXJCB\\$" + strconv.Itoa(qualification) + "\" value=\"(\\d+)\"")
		form["JiangXJCB$"+strconv.Itoa(qualification)] = []string{string(qualificationIdRegex.FindSubmatch(typePageBody)[1])}
	}
	_, _ = postFormWithUA(&client, "http://www.sz.shu.edu.cn/JiangXj/JiangXJStudentList.aspx", form)
	log.Println(student.Name, " apply for scholarship ", Types[scholarshipType])
}
